#!/bin/bash

curr_dir=`pwd`
goli_work_dir="/usr/local/sbin/goli"
go_download_link="https://go.dev/dl/go1.20.3.linux-amd64.tar.gz"
go_installer="go1.20.3.linux-amd64.tar.gz"

source_func() {
    source "/etc/profile"
}

exit_func() {
    exit 1
}

create_goli_user() {
    # Check if goli user exists
    if id "goli" &>/dev/null; then
        echo "User 'goli' already exists"
    else
        echo "Creating system user 'goli' ..."
        # Create system user with no login shell and no home directory
        useradd -r -s /usr/sbin/nologin -d /goli -c "Goli service user" goli
        if [ $? -eq 0 ]; then
            echo "User 'goli' created successfully"
        else
            echo "Failed to create user 'goli'"
            exit_func
        fi
    fi
}

install_go(){
    which go
    if [ $? -ne 0 ]; then 
        echo "Go does not exists!"
        echo "Downloading go ..."
    
        wget "${go_download_link}" -O "/tmp/${go_installer}"

        echo "Extraction ..."
        [ $? -eq 0 ] && rm -rf /usr/local/go && tar -C /usr/local -xzf "/tmp/${go_installer}"
        
        echo "Installing go ..."
        export PATH=$PATH:/usr/local/go/bin
        [ $? -eq 0 ] && grep "from-goli" "/etc/profile" 
        if [ $? -ne 0 ]; then
            echo 'export PATH=$PATH:/usr/local/go/bin # from-goli' >> "/etc/profile"
            echo 'export PATH=$PATH:/usr/local/go/bin # from-goli' >> "$HOME/.profile"
        fi

        /usr/local/go/bin/go version 2>/dev/null
        if [ $? -eq 0 ]; then  
            echo "Go successfully installed." 
        else
            echo "An error occured during the installation of 'Go' "
            exit_func 
        fi
    else
        echo "Go is already installed"
    fi
}

compile_and_install_binaries() {
    ## Compile and Install go binary
    cd "${curr_dir}/goli" && /usr/local/go/bin/go get && /usr/local/go/bin/go build -o "${goli_work_dir}/goli" main.go && cd -
}

install() {
    echo "========================================="
    echo "Installing Goli ..."
    echo "========================================="
    
    # Create goli system user first
    create_goli_user
    
    mkdir -p /usr/local/sbin/goli

    # Auto-generate auth key (no user interaction needed)
    default_authkey=`openssl rand -base64 64 | openssl sha256 | awk '{ print $2 }'`
    auth_key="${default_authkey}"

    # Generate one-time setup password (12 characters, alphanumeric)
    setup_password=`openssl rand -base64 12 | tr -d "=+/" | cut -c1-12`
    
    # Check if go is installed
    install_go

    echo "Compiling Goli binary ..."
    compile_and_install_binaries

    # Create Goli directories
    mkdir -p /goli/config
    mkdir -p /goli/data
    
    # Create Goli Toml config file with setup_complete flag set to false
    echo "s/dummy_key/${auth_key}/1" > "${curr_dir}/utils/rule_1.sed"
    sed -f "${curr_dir}/utils/rule_1.sed" "${curr_dir}/utils/config.toml" > "/goli/config/config.toml"
    
    # Add setup password to config file
    echo "setup_password = \"${setup_password}\"" >> "/goli/config/config.toml"

    # Set ownership of all goli directories and files to goli user
    echo "Setting ownership of Goli directories to goli user ..."
    chown -R goli:goli /goli
    chown -R goli:goli /usr/local/sbin/goli
    chmod 755 /goli
    chmod 755 /goli/config
    chmod 755 /goli/data
    chmod 644 /goli/config/config.toml
    chmod 755 /usr/local/sbin/goli
    chmod 755 /usr/local/sbin/goli/goli

    # Create Goli service file
    goli_work_dir_for_sed="\/usr\/local\/sbin\/goli"
    echo "s/work_dir/${goli_work_dir_for_sed}/1;s/exec_start/${goli_work_dir_for_sed}\/goli/1" > "${curr_dir}/utils/rule_2.sed"
    sed -f "${curr_dir}/utils/rule_2.sed" "${curr_dir}/utils/goli.service" > "/etc/systemd/system/goli.service"

    # Reload systemd and start service
    systemctl daemon-reload
    systemctl enable --now goli.service

    if [ $? -eq 0 ]; then
        echo ""
        echo "========================================="
        echo "Goli successfully installed!"
        echo "========================================="
        echo ""
        echo "Installation Summary:"
        echo "  - System user 'goli' created"
        echo "  - All Goli files owned by 'goli' user"
        echo "  - Service configured to run as 'goli' user"
        echo ""
        echo "Service Status:"
        systemctl status goli.service --no-pager -l | head -n 5
        echo ""
        echo "Next Steps:"
        echo "1. Access the Goli UI at: http://localhost:8125"
        echo "2. Complete the initial setup in the UI:"
        echo "   - Enter the setup password below"
        echo "   - Configure your admin user (default: 'goli')"
        echo "   - Update application settings"
        echo "   - Configure tool parameters"
        echo ""
        echo "========================================="
        echo "IMPORTANT: ONE-TIME SETUP PASSWORD"
        echo "========================================="
        echo ""
        echo "  ${setup_password}"
        echo ""
        echo "⚠️  SECURITY WARNING:"
        echo "   - This password is required to complete the initial setup"
        echo "   - It will be invalidated after successful setup"
        echo "   - Store it securely and do not share it"
        echo ""
        echo "Your temporary auth key (for API access):"
        echo "  ${auth_key}"
        echo ""
        echo "This key is stored in: /goli/config/config.toml"
        echo "After completing setup in the UI, you can manage users and settings from there."
        echo ""
        echo "Note: The Goli service runs as the 'goli' system user for security."
        echo "      All files in /goli/ are owned by this user."
        echo ""
    else
        echo "Goli installation failed. Cleaning up..."
        remove
    fi

    exit 0
}

update() {
    echo "Updating Goli ..."
    
    # Ensure goli user exists
    create_goli_user
    
    compile_and_install_binaries
    
    # Ensure ownership is correct after update
    chown -R goli:goli /usr/local/sbin/goli
    chmod 755 /usr/local/sbin/goli/goli
    
    systemctl restart goli.service
    echo "Goli successfully updated."
    exit 0
}

remove() {
    echo "Removing Goli ..."

    # Stop the service first
    systemctl stop goli.service 2>/dev/null
    
    # Remove the service from the system
    systemctl disable goli.service 2>/dev/null

    rm -f "/etc/systemd/system/goli.service"
    rm -rf "/goli/"
    rm -rf /usr/local/sbin/goli
    
    systemctl daemon-reload
    
    echo "Goli successfully removed."
    exit 0
}

source_func

if [ "$(whoami)" != "root" ]; then
    echo "Error: This script must be run as root"
    exit_func
fi

# Main menu
if [ ! -d "/usr/local/sbin/goli" ]; then
    echo -e "\n1. Install Goli\nq. Quit"
else
    echo -e "\n1. Update Goli\n2. Remove Goli\nq. Quit"
fi

read -p "Select an option: " to_do

case "${to_do}" in
    "1")
        [ -d "/usr/local/sbin/goli" ] && update || install
        ;;
    "2")
        [ -d "/usr/local/sbin/goli" ] && remove || (echo "Goli is not installed." && exit_func)
        ;;
    "q"|"Q")
        echo "Exiting..."
        exit 0
        ;;
    *)
        echo "Unknown option"
        exit_func
        ;;
esac

