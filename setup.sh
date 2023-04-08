#!/bin/bash

function source_func() {
    source "/etc/profile"
}

function exit_func() {
    exit 1
}

function install() {
    echo "Installing Goli ..."
    # 64 | openssl sha256 | awk '{ print $2 }'

    default_authkey=`openssl rand -base64 64 | openssl sha256 | awk '{ print $2 }'`

    read -p "Enter an authorization key (default '${default_authkey}' ): " AUTHKEY

    [ -z "$AUTHKEY" ] && auth_key="${default_authkey}" || auth_key="$AUTHKEY"

    # check if go is Intalled
    echo ""

    go_download_link="https://go.dev/dl/go1.20.3.linux-amd64.tar.gz"
    go_installer="go1.20.3.linux-amd64.tar.gz"

    if [ ! $(go version 2>/dev/null) ]; then 
        echo "Go does not exists!"
        echo "Downloading go ..."
    
        wget "${go_download_link}" -O "/tmp/${go_installer}"

        echo "Extraction ..."
        [ $? -eq 0 ] && rm -rf /usr/local/go && tar -C /usr/local -xzf "/tmp/${go_installer}"
        
        echo "Installing go ..."
        [ $? -eq 0 ] && echo 'export PATH=$PATH:/usr/local/go/bin' >> "/etc/profile"
        [ $? -eq 0 ] && source_func && echo "null" >/dev/null

        if [ ! $(go version 2>/dev/null) ]; then  
            echo "Go successfully installed." 
        else
            echo "An error occured during the installation of 'Go' "
            exit_func 
        fi
    else
        echo "Go is already installed"
    fi

   

    goli_work_dir="/usr/local/sbin/goli"
    mkdir -p /usr/local/sbin/goli

    ## Change port possibility


    ## Compile and Install go binary
    cd ./goli && go get && go build -o "${goli_work_dir}/goli" main.go && cd ..

    ## Create Goli Toml config file
    mkdir -p /goli/config
    echo "s/dummy_key/${auth_key}/1" > "./utils/rule_1.sed"
    sed -f "./utils/rule_1.sed" ./utils/config.toml > "/goli/config/config.toml"


    ## Create Goli service file
    goli_work_dir_for_sed="\/usr\/local\/sbin\/goli"
    echo "s/work_dir/${goli_work_dir_for_sed}/1;s/exec_start/${goli_work_dir_for_sed}\/goli/1" > "./utils/rule_2.sed"
    sed -f "./utils/rule_2.sed" ./utils/goli.service > "/etc/systemd/system/goli.service"

    systemctl enable --now goli.service

    echo "Goli Action helper successfully installed."

    exit 0
}

function remove() {
    echo "Removing Goli ..."

    # Stop the service first
     systemctl stop goli.service
    
    # remove the service from the system
    systemctl disable goli.service

    rm "/etc/systemd/system/goli.service" && rm -rf "/goli/" && rm -rf /usr/local/sbin/goli

    echo "Goli Action helper successfully removed."

    exit 0
}


source_func

if [ "$(whoami)" != "root" ]; then
    echo "Script must be run as user: root"
    exit_func
fi


[ ! -d "/usr/local/sbin/goli" ] && echo -e "1 Install Goli\nq Quit the program" || echo -e "1 Remove Goli\nq Quit the program"

read to_do

if [ "${to_do}" == "1" ]; then
    [ -d "/usr/local/sbin/goli" ] && remove || install

elif [ "${to_do}" == "q" ]; then
    exit_func
else 
    echo "Unknown option"
    exit_func
fi


