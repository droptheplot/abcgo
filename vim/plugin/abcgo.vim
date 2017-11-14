function! ABCGoBackground(channel, reports)
  if exists("g:abcgo_max")
    let abcgo_max = g:abcgo_max
  else
    let abcgo_max = 10
  endif

  let current_file = expand("%:p")

  sign define abcgo text=#! texthl=Visual
  execute "sign unplace * file=" . current_file

  for report in split(a:reports, "\n")
    let report = split(report, " ")
    if report[2] >= abcgo_max
      echo report
      execute ":sign place 63278 line=" . report[1] . " name=abcgo file=" . current_file
    endif
  endfor
endfunction

function! ABCGo()
  call job_start("go run main.go -path " . expand("%:p"), {'callback': 'ABCGoBackground', 'mode': 'raw'})
endfunction

augroup abcgo_autocmd
  autocmd BufWrite *.go call ABCGo()
  autocmd BufRead *.go call ABCGo()
augroup END
