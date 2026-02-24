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

## Release automatica no GitHub

Ao criar uma tag `v*`, o GitHub Actions gera e publica bundles Ubuntu (`amd64` e `arm64`) na Release.

Exemplo:

```bash
git tag v0.1.0
git push origin v0.1.0
```
