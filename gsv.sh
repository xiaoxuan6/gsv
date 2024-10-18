#!/bin/bash

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')

function install() {
  UNAME="${OS}_${ARCH}"
  URL=$(curl -s https://api.github.com/repos/xiaoxuan6/gsv/releases/latest | grep "browser_download_url" | grep "tar.gz" | cut -d '"' -f 4 | grep "$UNAME")
  if [ -z "$URL" ]; then
    echo "Unsupported platform: $(uname -s) $(uname -m)"
    exit 1
  fi

  curl -L -O "$URL"
  FILENAME=$(echo "$URL" | cut -d '/' -f 9)
  if [ ! -f "$FILENAME" ]; then
    echo "url: $URL"
    echo "filename $FILENAME dose not exist"
    exit 1
  fi

  tar xf "$FILENAME"
  rm "$FILENAME"

  chmod +x "gsv"
  mv "./gsv" "/usr/local/gsv"

  if [ ! -x "$(command -v gsv)" ]; then
    # shellcheck disable=SC2016
    if ! grep -q 'export PATH="\$PATH:/usr/local"' ~/.bashrc; then
        # shellcheck disable=SC2016
        echo 'export PATH="$PATH:/usr/local"' >> ~/.bashrc

        # shellcheck disable=SC1090
        source ~/.bashrc
    fi
  fi
  echo "gsv 安装成功！"
  echo
  echo -e "\e[31m使用 gsv 报错：Command 'gsv' not found\e[0m，请重新打开终端或者在终端直接执行：'source ~/.bashrc'"
}

function remove() {
  rm "/usr/local/gsv"
  echo "gsv 卸载成功！"
}

case $1 in
install)
  install
  ;;
remove)
  remove
  ;;
*)
  echo "Not found $1 option"
  echo "Usage: $0 {install|remove}"
  echo ""
  exit 1
  ;;
esac
