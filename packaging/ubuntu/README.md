# Stashclip (Ubuntu Bundle)

## Conteudo
- `bin/stashclip`: binario Linux.
- `install.sh`: instala binario, wrapper de popup e service user.
- `uninstall.sh`: remove instalacao.
- `stashclip.service`: service `systemd --user`.

## Instalacao rapida
```bash
tar -xzf stashclip-<versao>-ubuntu-<arch>.tar.gz
cd stashclip-<versao>-ubuntu-<arch>
./install.sh --install-deps
```

## Uso
- Daemon inicia automaticamente via `systemd --user`.
- Popup on-demand:
```bash
~/.local/bin/stashclip-popup
```

## Atalho global
Configure o atalho global para executar:
```bash
~/.local/bin/stashclip-popup
```
