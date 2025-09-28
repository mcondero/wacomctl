# wacomctl

Um utilitário de linha de comando para controlar tablets Wacom no Linux.

## Descrição

O `wacomctl` é uma ferramenta simples escrita em Go que permite mapear e controlar a stylus de tablets Wacom para diferentes monitores ou desabilitá-la/habilitá-la conforme necessário.

## Pré-requisitos

- Sistema Linux com X11
- Tablet Wacom conectado
- Pacotes `xsetwacom` e `xinput` instalados
- Go 1.x (para compilação)

## Instalação

1. Clone ou baixe o código fonte
2. Compile o executável:

   ```bash
   go build -o wacomctl wacomctl.go
   ```

3. Torne o executável disponível no PATH (opcional):

   ```bash
   sudo mv wacomctl /usr/local/bin/
   ```

## Uso

```bash
wacomctl [vga|hdmi|both|off|on]
```

### Comandos

- **`vga`** - Mapeia a stylus para o monitor VGA conectado (ex: VGA-1)
- **`hdmi`** - Mapeia a stylus para o monitor HDMI conectado (ex: HDMI-1)
- **`both`** - Mapeia a stylus para todos os monitores ativos (desktop completo)
- **`off`** - Desliga a stylus (desabilita o dispositivo)
- **`on`** - Liga a stylus (habilita o dispositivo)

### Exemplos

```bash
# Mapear stylus para monitor VGA
wacomctl vga

# Mapear stylus para monitor HDMI
wacomctl hdmi

# Mapear stylus para todos os monitores
wacomctl both

# Desligar a stylus
wacomctl off

# Ligar a stylus
wacomctl on
```

## Funcionamento

O programa utiliza as seguintes ferramentas do sistema:

- **`xsetwacom`** - Para detectar e configurar dispositivos Wacom
- **`xrandr`** - Para listar monitores conectados
- **`xinput`** - Para habilitar/desabilitar o dispositivo

### Detecção automática

- Detecta automaticamente o ID do dispositivo stylus através do comando `xsetwacom --list devices`
- Identifica monitores VGA e HDMI conectados através do comando `xrandr --listmonitors`
- Mapeia a área ativa da stylus para o monitor especificado

## Limitações

- Funciona apenas em sistemas X11 (não Wayland)
- Requer que os monitores estejam conectados e ativos
- Detecta apenas o primeiro dispositivo stylus encontrado
- Suporta apenas monitores VGA e HDMI (não DisplayPort, DVI, etc.)

## Solução de problemas

### Dispositivo não encontrado

```text
Dispositivo de stylus não encontrado.
```

- Verifique se o tablet Wacom está conectado
- Execute `xsetwacom --list devices` para verificar se o dispositivo é reconhecido

### Monitor não encontrado

```text
Nenhum monitor VGA/HDMI encontrado.
```

- Verifique se o monitor está conectado e ativo
- Execute `xrandr --listmonitors` para ver os monitores disponíveis

### Erro de mapeamento

```text
Erro ao mapear stylus: ...
```

- Verifique se você tem permissões adequadas
- Certifique-se de que o `xsetwacom` está instalado

## Licença

Este projeto é de domínio público. Use livremente.
