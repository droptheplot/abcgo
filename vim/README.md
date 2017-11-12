# ABCGo

Vim plugin for ABCGo.

## Getting Started

### Installation

#### Vundle

Place this in your `.vimrc`:

```vim
Plugin 'droptheplot/abcgo', {'rtp': 'vim/'}
```

Run the following in Vim:

```vim
:source %
:PluginInstall
```

#### VimPlug

Place this in your `.vimrc`:

```vim
Plug 'droptheplot/abcgo', { 'rtp': 'vim/' }
```

Run the following in Vim:

```vim
:source %
:PlugInstall
```

### Configuration

Change maximum allowed ABC score:

```vim
let g:abcgo_max = 20
```
