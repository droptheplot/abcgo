function! ABCGoBackground(channel, reports)
  if !exists("g:abcgo_max")
    let g:abcgo_max = 20
  endif

  if !exists("b:abcgo_count")
    let b:abcgo_count = 0
  endif

  let current_file = expand("%:p")
  sign define abcgo text=c texthl=WarningMsg

  if strlen(a:reports) > 0
    let i = 0

    while i <= b:abcgo_count
      let i += 1

      execute "sign unplace 9314 file=" . current_file
    endwhile
  endif

  let b:abcgo_count = 0

  for report in split(a:reports, "\n")
    let report = split(report, "\t")

    if report[3] >= g:abcgo_max
      let b:abcgo_count += 1
      execute ":sign place 9314 line=" . report[1] . " name=abcgo file=" . current_file
    endif
  endfor
endfunction

function! ABCGo()
  call job_start($GOPATH . "/bin/abcgo -format raw -path " . expand("%:p"), {'callback': 'ABCGoBackground', 'mode': 'raw'})
endfunction

augroup abcgo_autocmd
  autocmd!

  autocmd BufWritePost *.go call ABCGo()
  autocmd BufReadPost *.go call ABCGo()
augroup END
