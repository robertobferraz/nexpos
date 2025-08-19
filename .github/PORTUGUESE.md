<!-- BANNER -->
<p align="center">
  <img src=".github/banner.png" alt="Nexpos Banner" width="100%" />
</p>

<!-- BADGES -->
<p align="center">
  <img src="https://img.shields.io/badge/language-Go-blue?logo=go" />
  <img src="https://img.shields.io/badge/docker-ready-blue?logo=docker" />
  <img src="https://img.shields.io/badge/license-MIT-green?logo=open-source-initiative" />
  <img src="https://img.shields.io/badge/status-active-success?logo=github" />
</p>

---

# 📘 Nexpos

> Um serviço monolítico para gerenciar um sistema de ponto de venda (POS) online, desenvolvido em Go com integração a serviços externos (Firebase), banco de dados PostgreSQL, cache Redis e Swagger para documentação de API.
---

## 📑 Sumário

- [Introdução](#-introdução)
- [Tecnologias](#-tecnologias)
- [Arquitetura](#-arquitetura)
- [Instalação](#-instalação)
- [Configuração](#-configuração)
- [Uso](#-uso)
- [Contribuição](#-contribuição)
- [Licença](#-licença)

---

## 📖 Introdução

**NexPOS** é um sistema de ponto de venda (PDV) de código aberto projetado para gerenciar vendas online, estoque e interações com clientes.
Desenvolvido com foco em **escalabilidade**, **segurança** e **práticas de arquitetura limpa**, ele fornece uma base sólida para soluções de e-commerce. Seja você um pequeno varejista ou um grande varejista, o NexPOS oferece ferramentas para processamento de pedidos, gerenciamento de usuários e relatórios em tempo real.
---

## ✨ Recursos

- Suporte multiusuário com autenticação Firebase
- Gerenciamento de inventário em tempo real com cache Redis
- API RESTful com documentação Swagger
- Integração de processamento de pedidos e pagamentos
- Relatórios e análises personalizáveis

---

## 🛠 Tecnologias

- [Go](https://go.dev/) – Linguagem de programação principal
- [Docker](https://www.docker.com/) – Conteinerização
- [GORM](https://gorm.io/) – ORM para PostgreSQL
- [PostgreSQL](https://www.postgresql.org/) – Banco de dados relacional
- [Redis](https://redis.io/) – Armazenamento de dados em memória
- [Firebase](https://firebase.google.com/) – Autenticação e serviços externos
- [Swagger](https://swagger.io/) – Documentação da API
- **Arquitetura Limpa / Ports and Adapters**
---

## 🏗 Arquitetura

O projeto segue os princípios da **Arquitetura Limpa**:

- **Domínio** → Regras de Negócio
- **Casos de Uso** → Casos de Uso da Aplicação
- **Adaptadores** → Interfaces Externas (BD, APIs)
- **Infraestrutura** → Configurações e Frameworks

---

## ⚙️ Instalação

```bash
# Clone o repositório
git clone https://github.com/robertobferraz/nexpos.git

# Acesse o diretório
cd nexpos

# Inicie com Docker
make up 
```

---

## 🔧 Configuração

As variáveis de ambiente estão no arquivo `.env`.  
Exemplo:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=secret
DB_NAME=nexpos
```

---

## 🚀 Uso

```bash
# Rodar aplicação localmente
go run main.go
```

Acesse: [http://localhost:8080](http://localhost:8080)

---

## 🤝 Contribuição

Contribuições são bem-vindas!  
Siga estes passos:

1. Faça um fork do projeto
2. Crie uma branch (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova feature'`)
4. Faça push (`git push origin feature/nova-feature`)
5. Abra um Pull Request

---

## 📜 Licença

Este projeto está sob a licença [MIT](LICENSE).

---