# unix-socket-path-limit

> A reproducible proof of the AF_UNIX socket path length limit — and the
> runtime-directory fix — distilled from [kubernetes/minikube#21087][issue].

**Languages:** [English](#english) · [Português](#português)

---

## English

A minimal, reproducible demonstration of the **AF_UNIX socket path length
limit** (`sun_path`) and the runtime-directory pattern that works around it.

### The problem

`struct sockaddr_un.sun_path` is a fixed-size array — **108 bytes on Linux**,
**104 on macOS/BSD**. It is a hard kernel ABI limit: a socket whose path
exceeds it cannot be bound.

This is the root cause of [minikube#21087][issue] — per-machine QMP sockets
stored under a long `$MINIKUBE_HOME` blew the limit, and `qemu` refused to
start.

### How the fix works

`runtimedir.Resolve` builds each per-machine socket path under a short base
directory:

```
<base>/minikube/<hash>/<socket-name>
```

`<hash>` is the first 32 hex characters of `sha256(machineName)`. It
compresses an arbitrary-length machine name into a fixed-length directory
component, so the resolved path stays under `sun_path` no matter how long
the input name is. The platform base is chosen at runtime — no build tags.

### Layout

| Path                     | Contents                                                |
| ------------------------ | ------------------------------------------------------- |
| `runtimedir/`            | the helper (the fix) + table-driven unit tests          |
| `integration/`           | Testcontainers tests — real `qemu` in a Linux container |
| `docker/Dockerfile.qemu` | the container image (Debian + `qemu-system`)            |

### Run it

| Command                 | What it does                              |
| ----------------------- | ----------------------------------------- |
| `make test`             | fast unit tests — no Docker               |
| `make test-integration` | `qemu` in a real container — needs Docker |
| `make verify`           | both, in order                            |

`make test-integration` reproduces the bug — `qemu` rejects an over-limit
socket path — and proves the fix — `qemu` accepts the path that `runtimedir`
resolves.

### Not in scope

This is a teaching PoC, not a library. The production minikube helper also
handles the Linux `$XDG_RUNTIME_DIR` branch, a Windows branch, lifecycle
cleanup, and a pure-vs-side-effecting API split. See [#21087][issue] for the
real code.

---

## Português

Uma demonstração mínima e reproduzível do **limite de comprimento de caminho
do socket AF_UNIX** (`sun_path`) e do padrão de diretório de runtime que
contorna o problema.

### O problema

`struct sockaddr_un.sun_path` é um array de tamanho fixo — **108 bytes no
Linux**, **104 no macOS/BSD**. É um limite rígido da ABI do kernel: um socket
cujo caminho excede esse limite não pode ser vinculado (`bind`).

Essa é a causa raiz da [minikube#21087][issue] — sockets QMP por máquina
guardados sob um `$MINIKUBE_HOME` longo estouravam o limite, e o `qemu` se
recusava a iniciar.

### Como a correção funciona

`runtimedir.Resolve` constrói o caminho de cada socket por máquina sob um
diretório base curto:

```
<base>/minikube/<hash>/<socket-name>
```

`<hash>` são os primeiros 32 caracteres hexadecimais de
`sha256(machineName)`. Ele comprime um nome de máquina de comprimento
arbitrário em um componente de diretório de tamanho fixo, então o caminho
resolvido permanece dentro do `sun_path` por mais longo que seja o nome de
entrada. A base por plataforma é escolhida em tempo de execução — sem build
tags.

### Estrutura

| Caminho                  | Conteúdo                                                   |
| ------------------------ | ---------------------------------------------------------- |
| `runtimedir/`            | o helper (a correção) + testes unitários table-driven      |
| `integration/`           | testes com Testcontainers — `qemu` real em container Linux |
| `docker/Dockerfile.qemu` | a imagem do container (Debian + `qemu-system`)             |

### Como executar

| Comando                 | O que faz                                    |
| ----------------------- | -------------------------------------------- |
| `make test`             | testes unitários rápidos — sem Docker        |
| `make test-integration` | `qemu` em container real — precisa de Docker |
| `make verify`           | os dois, em ordem                            |

`make test-integration` reproduz o bug — o `qemu` rejeita um caminho de
socket acima do limite — e prova a correção — o `qemu` aceita o caminho que
o `runtimedir` resolve.

### Fora de escopo

Este é um PoC didático, não uma biblioteca. O helper de produção do minikube
também trata o branch `$XDG_RUNTIME_DIR` no Linux, um branch para Windows,
limpeza de ciclo de vida e a separação de API pura vs. com efeitos
colaterais. Veja a [#21087][issue] para o código real.

[issue]: https://github.com/kubernetes/minikube/issues/21087
