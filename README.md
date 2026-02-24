# Stashclip

Stash it. Find it. Paste it.

## Ubuntu Bundle

Build do pacote Ubuntu (amd64):

```bash
./scripts/build_ubuntu_bundle.sh amd64
```

Arquivos gerados em `dist/`:
- `stashclip-<versao>-ubuntu-amd64.tar.gz`
- `stashclip-<versao>-ubuntu-amd64.tar.gz.sha256`

Instalacao no Ubuntu:

```bash
tar -xzf stashclip-<versao>-ubuntu-amd64.tar.gz
cd stashclip-<versao>-ubuntu-amd64
./install.sh --install-deps
```
