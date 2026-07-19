Para publicar o **Shelloma** nas lojas de aplicativos do **Ubuntu (Snap Store)** e do **Flatpak (Flathub)**, o seu repositório já conta com quase toda a estrutura necessária (manifestos `.yml`, ícones, desktop e metadados AppStream).

Abaixo está o passo a passo detalhado do que é necessário para realizar a publicação em cada uma das lojas:

---

## 1. 🛍️ Ubuntu Snap Store (Canonical Snapcraft)

A **Snap Store** da Canonical é a loja padrão do Ubuntu. Uma vez publicado lá, os usuários podem instalar seu aplicativo digitando `sudo snap install shelloma` ou buscando diretamente no **Ubuntu Software / Snap Store**.

### O que você precisa fazer:

1. **Criar uma conta no Canonical SSO**:
   - Acesse o portal oficial [dashboard.snapcraft.io](https://dashboard.snapcraft.io/) ou [snapcraft.io](https://snapcraft.io/) e crie uma conta gratuita.

2. **Reservar o Nome da Aplicação**:
   - No seu terminal (ou no portal web do Snapcraft), faça login e reserve o nome `shelloma`:
     ```bash
     snapcraft login
     snapcraft register shelloma
     ```

3. **Gerar Credenciais para o GitHub Actions**:
   - Exporte um token de acesso para o pipeline automatizado:
     ```bash
     snapcraft export-login --channels=stable snapcraft.login
     ```
   - Copie o conteúdo gerado no arquivo `snapcraft.login` e adicione como um **GitHub Secret** no seu repositório:
     - Acesse no GitHub: `Settings > Secrets and variables > Actions > New repository secret`.
     - Nome: `SNAPCRAFT_STORE_CREDENTIALS`.
     - Valor: Cole o texto do arquivo `snapcraft.login`.

4. **Habilitar a Publicação Automática no GoReleaser**:
   - No arquivo [.goreleaser.yaml](file:///home/moises/Projetos/shelloma/.goreleaser.yaml), altere a opção de publicação de `publish: false` para `publish: true`:
     ```yaml
     snapcrafts:
       - id: shelloma
         name: shelloma
         # ...
         publish: true
     ```
   - Passe o segredo `SNAPCRAFT_STORE_CREDENTIALS` no passo do GoReleaser em [.github/workflows/release.yml](file:///home/moises/Projetos/shelloma/.github/workflows/release.yml).

---

## 2. 📦 Flathub (Loja Oficial do Flatpak)

O **Flathub** ([flathub.org](https://flathub.org)) é o repositório e loja central de aplicativos Flatpak, pré-instalado ou facilmente habilitado em praticamente todas as distribuições Linux (Fedora, Arch Linux, Linux Mint, Debian, Ubuntu, SteamOS, etc.).

### O que você precisa fazer:

1. **Estrutura do Repositório (Já Concluída ✅)**:
   - Já possuímos o manifesto em [scripts/org.shelloma.Shelloma.yml](file:///home/moises/Projetos/shelloma/scripts/org.shelloma.Shelloma.yml).
   - Já possuímos o arquivo AppStream em [scripts/org.shelloma.Shelloma.appdata.xml](file:///home/moises/Projetos/shelloma/scripts/org.shelloma.Shelloma.appdata.xml).
   - Já possuímos o atalho desktop em [scripts/shelloma.desktop](file:///home/moises/Projetos/shelloma/scripts/shelloma.desktop).

2. **Submeter o App via Pull Request no Flathub**:
   - Faça um Fork do repositório de submissões do Flathub no GitHub: **[flathub/flathub](https://github.com/flathub/flathub)**.
   - Crie uma branch nomeada com o seu App ID: `new-pr-org.shelloma.Shelloma`.
   - Adicione o manifesto `org.shelloma.Shelloma.yml` (apontando as fontes para as releases do seu GitHub).
   - Abra uma **Pull Request** para a equipe do Flathub.

3. **Revisão e Publicação**:
   - O bot do Flathub e a comunidade revisarão os metadados e permissões.
   - Uma vez aprovado, o Flathub criará um repositório dedicado em `github.com/flathub/org.shelloma.Shelloma`.
   - O aplicativo estará disponível para qualquer usuário Linux no mundo via:
     ```bash
     flatpak install flathub org.shelloma.Shelloma
     ```

---

### 💡 Resumo do Próximo Passo

Se você desejar publicar nas lojas oficiais:
- Para a **Snap Store**: Crie a conta no Snapcraft, registre o nome `shelloma` e gere a chave `SNAPCRAFT_STORE_CREDENTIALS`.
- Para o **Flathub**: Nós podemos preparar a Pull Request de submissão do manifesto Flatpak.