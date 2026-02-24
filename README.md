# Stashclip

Clipboard manager para Linux (X11 e Wayland) com daemon em background e popup on-demand para escolher e copiar itens salvos.

## O que o projeto faz

- Captura automaticamente mudanças do clipboard e salva em storage local.
- Mantém um daemon rodando em segundo plano.
- Abre popup de seleção para copiar novamente qualquer item salvo.
- Funciona com provedores gráficos comuns (`yad`, `zenity`, `kdialog`).

## Download e instalação (Ubuntu)

Baixe o pacote da release no GitHub (`stashclip-<versao>-ubuntu-<arch>.tar.gz`) e rode:

```bash
tar -xzf stashclip-<versao>-ubuntu-amd64.tar.gz
cd stashclip-<versao>-ubuntu-amd64
./install.sh --install-deps
```

Depois disso:
- Binário: `~/.local/bin/stashclip`
- Comando para atalho global: `~/.local/bin/stashclip-popup`

## Comandos principais

```bash
stashclip daemon start
stashclip daemon status
stashclip daemon stop

stashclip list
stashclip pick
stashclip pick 3
stashclip popup
stashclip clear
```

## Build local do bundle Ubuntu

```bash
./scripts/build_ubuntu_bundle.sh amd64
./scripts/build_ubuntu_bundle.sh arm64
```

Arquivos gerados em `dist/`:
- `stashclip-<versao>-ubuntu-amd64.tar.gz`
- `stashclip-<versao>-ubuntu-arm64.tar.gz`
- `*.sha256`

## Release automática no GitHub

Ao criar uma tag `v*`, o GitHub Actions gera e publica os bundles Ubuntu na Release.

```bash
git tag v0.1.0
git push origin v0.1.0
```
