function! s:ABCGo()
  if exists("g:abcgo_max")
    let abcgo_max = g:abcgo_max
  else
    let abcgo_max = 10
  endif

  let current_file = expand("%:p")
  let reports = system("go run main.go -path " . current_file)

  sign define abcgo text=#! texthl=Visual

  execute "sign unplace * file=" . current_file

  for report in split(reports, "\n")
    let report = split(report, " ")
    if report[2] >= abcgo_max
      execute ":sign place 1 line=" . report[1] . " name=abcgo file=" . current_file
    endif
  endfor
endfunction
command! ABCGo call s:ABCGo()

augroup abcgo_autocmd
  autocmd BufWrite *.go ABCGo
  autocmd BufRead *.go ABCGo
augroup END
