# 🐚 Shelloma

> **Tradutor de Linguagem Natural para Comandos Shell do Linux (Debian/Ubuntu) alimentado por Ollama local.**

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Linux%20(Debian%2FUbuntu)-E6007E?logo=debian)](https://www.debian.org)
[![Ollama](https://img.shields.io/badge/LLM-Ollama-black?logo=ollama)](https://ollama.com)

**[English](README.md) | [Português](README_pt.md) | [Español](README_es.md)**

O **Shelloma** é uma aplicação CLI nativa e ultrarrápida desenvolvida em Go que traduz suas instruções em linguagem natural diretamente para comandos de terminal (Bash/Zsh) executáveis no Linux, utilizando o seu próprio modelo local via **Ollama**.

---

## 🛠️ Como Funciona?

1. **Interpretação**: O Shelloma captura sua instrução em português (ou outro idioma configurado), detecta automaticamente as especificidades do seu sistema operacional Linux (distribuição, versão, usuário atual, diretório de trabalho) e envia o contexto para a API do Ollama local.
2. **Geração de Comando**: O modelo gera estritamente o comando Bash/Zsh correspondente.
3. **Menu Interativo**: O Shelloma exibe o comando formatado em um card visual e permite que você escolha se deseja **Executar**, **Pedir Explicação do Comando**, **Modificar**, **Copiar para Clipboard** ou **Sair**.
4. **Análise de Execução e Recuperação de Erros**: Após a execução, o Shelloma analisa o retorno do terminal. Se houver falha, ele identifica a causa e **sugere automaticamente um comando de correção**, retornando ao comando anterior de forma encadeada assim que o erro for corrigido.

---

## ✨ Recursos

- 🚀 **Nativo e Leve**: Executável único compilado em Go, sem necessidade de Python, Node.js ou runtime externo.
- 🔒 **100% Privado e Offline**: Nenhum dado ou comando sai da sua máquina.
- 🌐 **Internacionalização Nativa (i18n)**: Suporte completo para **Português Brasileiro (`pt`)**, **Inglês (`en`)** e **Espanhol (`es`)** através de arquivos JSON embarcados.
- 🤖 **Seleção Automática de Modelo**: Detecta os modelos instalados no Ollama e seleciona automaticamente o melhor modelo focado em código/shell disponível.
- 💡 **Detecção e Inicialização do Ollama**: Se o serviço Ollama estiver parado, o Shelloma avisa e oferece um comando rápido interativo para iniciá-lo (`sudo systemctl start ollama`) com retry loop automático.
- 📦 **Instalação Simplificada `.deb`**: Pacote Debian pronto com assistente interativo de configuração de idioma na instalação.

---

## ⚙️ Dependências

1. **Sistema Operacional**: Linux (Debian, Ubuntu ou distribuições derivadas).
2. **Ollama**: Serviço do Ollama instalado e ativo localmente.
3. **Go 1.20+** *(opcional)*: Apenas caso deseje compilar a aplicação a partir do código-fonte.

---

## 🦙 Como Instalar e Configurar o Ollama

O **[Ollama](https://ollama.com)** é a ferramenta aberta responsável por rodar modelos de inteligência artificial localmente no seu computador.

- **Website Oficial**: [https://ollama.com](https://ollama.com)
- **Repositório GitHub**: [https://github.com/ollama/ollama](https://github.com/ollama/ollama)

### Passo 1: Instalar o Ollama no Linux

No seu terminal Linux, execute o comando oficial de instalação:

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

### Passo 2: Verificar a Execução do Serviço

Geralmente o Ollama inicia automaticamente como serviço do sistema. Você pode verificar o status com:

```bash
sudo systemctl status ollama
```

*Se o serviço não estiver rodando, você pode iniciá-lo com `sudo systemctl start ollama` ou executar manualmente `ollama serve`.*

---

## 🎯 Modelos Recomendados para Shell e Código

Para obter as melhores respostas de comandos de terminal Linux, recomendamos utilizar modelos treinados para código. Execute o comando `ollama pull <modelo>` para baixar o modelo desejado:

1. **Qwen 2.5 Coder 1.5B** *(Altamente Recomendado - Leve e Ultrarrápido)*:
   ```bash
   ollama pull qwen2.5-coder:1.5b
   ```
2. **Qwen 2.5 Coder 7B** *(Excelente precisão para tarefas complexas)*:
   ```bash
   ollama pull qwen2.5-coder:7b
   ```
3. **DeepSeek Coder 6.7B**:
   ```bash
   ollama pull deepseek-coder:6.7b
   ```
4. **Llama 3.2 3B**:
   ```bash
   ollama pull llama3.2:3b
   ```

---

## 📦 Guia de Instalação e Downloads

Todos os binários e pacotes pré-compilados do Shelloma estão disponíveis para download na página oficial de **[Releases no GitHub](https://github.com/moisesfilho/shelloma/releases)** (`.deb`, `.rpm`, `.tar.gz`, `.snap`, `AppImage`).

### Opção 1: Download de Binários e Pacotes Pré-compilados (Recomendado)

Acesse a página de **[GitHub Releases](https://github.com/moisesfilho/shelloma/releases)** e escolha a opção ideal para a sua distribuição Linux:

- **Debian / Ubuntu / Linux Mint / Pop!_OS (`.deb`)**:
  ```bash
  sudo dpkg -i shelloma_1.0.0_amd64.deb
  ```
  *(Durante a instalação do `.deb`, será exibido um assistente no terminal para selecionar o idioma padrão).*

- **Fedora / RedHat / CentOS / RHEL (`.rpm`)**:
  ```bash
  sudo rpm -i shelloma_1.0.0_amd64.rpm
  ```

- **AppImage (Universal - Executável portável sem necessidade de instalação)**:
  ```bash
  chmod +x Shelloma-x86_64.AppImage
  ./Shelloma-x86_64.AppImage "listar arquivos do diretorio"
  ```

- **Arquivo Comprimido Tarball (`.tar.gz`)**:
  Baixe o `.tar.gz` (`amd64` ou `arm64`), extraia o binário e mova para seu `PATH` (ex: `~/.local/bin/`).

### Opção 2: Compilando a partir do Código-Fonte

```bash
# 1. Clonar o repositório
git clone https://github.com/moisesfilho/shelloma.git
cd shelloma

# 2. Compilar e instalar no diretório do usuário (~/.local/bin)
make build
make install-user
```

---

## 🚀 Guia de Uso Simplificado

### Uso Básico

Basta chamar o `shelloma` seguido da sua instrução entre aspas:

```bash
# Exemplos em Português
shelloma "listar todos os arquivos pdf da pasta downloads"
shelloma "verificar espaço em disco e uso de memória"
shelloma "criar uma pasta chamada fotos e mover todas as imagens png para ela"

# Exemplos em outros idiomas
shelloma -l en "list all active docker containers"
shelloma -l es "mostrar procesos que consumen mas cpu"
```

### Menu de Opções da CLI

Ao gerar o comando, o Shelloma exibirá o card e aguardará sua ação:

```text
┌────────────────────────────────────────────┐
│  ls -la ~/Downloads/*.pdf                 │
└────────────────────────────────────────────┘

Opções: [Enter/y: Executar] [e: Explicar] [m: Modificar] [c: Copiar] [q/n: Sair]:
```

- **Pressionar `Enter` ou `y`**: Executa o comando diretamente no terminal.
- **Digite `e`**: Solicita uma explicação linha por linha do comando ao Ollama.
- **Digite `m`**: Abre um prompt para você editar o comando antes de executar.
- **Digite `c`**: Copia o comando diretamente para a área de transferência do sistema.
- **Digite `q` ou `n`**: Cancela a operação de forma segura.

### Comandos de Configuração e Utilitários

```bash
# Exibir as configurações atuais do Shelloma
shelloma config

# Alterar o idioma padrão salvo (pt, en, es)
shelloma config set lang pt

# Definir um modelo específico do Ollama
shelloma config set model qwen2.5-coder:1.5b

# Listar os modelos do Ollama instalados na sua máquina
shelloma models
```

### Flags Disponíveis

```text
Opções:
  -l, --lang string    Idioma: en, pt, es (padrão: en)
  -m, --model string   Modelo Ollama a utilizar (ex: qwen2.5-coder:1.5b)
  -url string          URL da API do Ollama (padrão: http://localhost:11434)
  -y, --yes            Executar o comando gerado automaticamente sem confirmação
  -v, --version        Exibir versão do Shelloma
```

---

## 📄 Licença

Este projeto está licenciado sob a licença **MIT** - consulte o arquivo [LICENSE](LICENSE) para obter mais detalhes.
