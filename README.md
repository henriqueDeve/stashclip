# Stashclip

Clipboard manager para Linux (X11 e Wayland) com daemon em background e popup on-demand para escolher e copiar itens salvos.

## O que o projeto faz

- Captura automaticamente mudanças do clipboard e salva em storage local.
- Mantém um daemon rodando em segundo plano.
- Abre popup de seleção para copiar novamente qualquer item salvo.
- Funciona com provedores gráficos comuns (`yad`, `zenity`, `kdialog`).

## Instalação (Ubuntu)

### Opção 1: Pacote `.deb` (recomendado)

Baixe o arquivo da Release:
- `stashclip_<versao>_amd64.deb` (ou `arm64`)

Instale com:

```bash
sudo apt-get update
sudo apt-get install -y ./stashclip_<versao>_amd64.deb
systemctl --user daemon-reload
systemctl --user enable --now stashclip.service
```

### Opção 2: Instalação em 1 comando (`curl | bash`)

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/henrique/stashclip/main/scripts/install.sh)
```

Para instalar uma versão específica:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/henrique/stashclip/main/scripts/install.sh) v0.1.0
```

Depois da instalação:
- Binário: `/usr/bin/stashclip`
- Popup: `/usr/bin/stashclip-popup`

## Como usar (usuario final)

1. Instale o app.
2. Configure o atalho global `Ctrl+Alt+A` para executar:

```bash
/usr/bin/stashclip-popup
```

3. Sempre que quiser reutilizar um texto copiado, pressione `Ctrl+Alt+A`,
   escolha no popup e cole com `Ctrl+V`.

## Build local do bundle Ubuntu

```bash
./scripts/build_ubuntu_bundle.sh amd64
./scripts/build_ubuntu_bundle.sh arm64
./scripts/build_ubuntu_deb.sh amd64
./scripts/build_ubuntu_deb.sh arm64
```

Arquivos gerados em `dist/`:
- `stashclip-<versao>-ubuntu-amd64.tar.gz`
- `stashclip-<versao>-ubuntu-arm64.tar.gz`
- `stashclip_<versao>_amd64.deb`
- `stashclip_<versao>_arm64.deb`
- `*.sha256`

## Release automática no GitHub

Ao criar uma tag `v*`, o GitHub Actions gera e publica os bundles Ubuntu na Release.

```bash
git tag v0.1.0
git push origin v0.1.0
```
