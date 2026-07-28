export default {
    common: {
        login: "Login",
        logout: "Logout",
        loading: "Loading...",
        theme: "Theme",
        language: "Language",
        search: "Search",
        actions: "Actions",
        noData: "No data",
        cancel: "Cancel",
        create: "Create",
        created: "Created"
    },

    auth: {
        username: "Username",
        password: "Password",
        loginTitle: "Sign In",
        loginError: "Unable to sign in",
    },

    menu: {
        dashboard: "Dashboard",
        users: "Users",
        endpoints: "Endpoints",
        groups: "Groups",
        configtypes: "Configuration Types",
        configs: "Configurations",
        settings: "Settings"
    },

    pages: {
        dashboard: {
            description: "Welcome to the dashboard",
            quickActions: "Quick actions",
            systemStatus: "System status",
            recentActivity: "Recent activity",
        },
        users: {
            description: "Manage users",
            create: "Create user",
            username: "Username",
            group: "Group",
            role: "Role",
            configs: "Configs",
            limit: "Limit",
            created: "Создан",
            admin: "Administrator",
            user: "User",
            delete: "Delete user",

        },

        endpoints: {
            description: "Manage endpoints",
            create: "Create endpoint",
            name: "Name",
            address: "Address",
            port: "Port",
            groups: "Groups",
        },

        groups: {
            description: "Manage groups",
            create: "Create group",
            name: "Name group",
            groups_descriptions: "Description",
            users: "User",
            endpoints: "Endpoints",
        }, 

        configtypes: {
            name: "Name",
            description: "Manage configuration types",
            description_type: "Description",
            createbutton: "Create configuration type",
            edit: "Edit configuration type",
            create: "Create configuration type",
            active: "Status"
        },
        configtypes_create: {
            name: "Name",
            description: "Description",
            active: "Status",
            template: "Configuration template",
            instruction: "Instructions for use",
            links_clients: "Links apps",
            create: "Create configuration type",
        },
        configs: {
            description: "Manage configurations",
            create: "Create configuration",
            download: "Download",
            qrcode: "QR Code",
            view: "View",
            total: "Total",
            used: "Used",
            available: "Available",
        },
        configs_edit: {
            endpoint: "Endpoint",
            created: "Created",
            configs: "Configuration content",
            copy: "Copy",
            download: "Download",
            detail: "Connection details",
            address: "address",
            instruction: "Instructions for use",
            clients: "App links",
            QR_Code: "Scan QR Code in App"
        },
    },
    dialog: {

        users: {
            username: {
                value: "Username",
                description: " Min 3 symbols, just latin word, цифры и _ "
            },
            password: {
                value: "Password",
                description: "Min 3 symbols"
            },
            group: {
                value: "Group",
                description: "The group identifies available connection points for users"
            },
            role: "Role",
            configs: "Configurations",
            limit: "Limit",
            limit_configurations: "Limit Confugurations",
            limit_configurations_hint: "Maximum of configrutions for user",
            created: "Created",
            admin: {
                    value: "Administrator",
                    description: "Full access managment of system"
                },
            user: {
                    value: "User",
                    description: "Creation and management of various modules"
            },
            create: "Create user",
            info_title: "Information",
            info_hint: "Roles of users",
            info2_title: "Limits of configurations",
            info2_hint: "Default: 5 configurations",
            info2_hint2: "Recommended: 10-20 for active users",
            info2_hint3: "Max: 100 for trusted users",
            message_delete: "Are you sure delete user",
            cant_message_delete: "This action can't cancel"       
        },
        endpoints: {
            name: "Name",
            address: "Address",
            port: "Port",
            configuration_types: "Configuration Types",        
            groups: "Groups",
            name_dialog: "Create endpoint",
            examples_title: "Address Examples",
            config_types_title: "Configuration Types",
            config_types_description: "Select which configuration types will be available for this endpoint.",
            address_hint: "IP address or domain name",
            port_hint: "Default WARP port: 2408",
            config_types_hint: "Available configuration types for this endpoint (Ctrl+Click to select multiple)",
            groups_hint: "User groups that can access this endpoint. If no groups are selected, the endpoint will be available to all users.",
            name_placeholder: "Endpoint name",
            address_placeholder: "Host or IP address",
            port_placeholder: "Port",
            config_types_placeholder: "Select configuration types",
            create: "Create Endpoint",
        },
        groups: {
            name_dialog: "Create group",
            name: "Name",
            name_placeholder: "Group name",
            name_hint: "Uniq name group",
            description: "Description",
            description_placeholder: "Description group",
            description_hint: "Optional description group",
            users: "Users",
            endpoints: "Available endpoints",
            endpoints_hint: "Pick endpoints",
            create: "Create configuration",
            
        },
        config: {
            name_dialog: "Create new configuration",
            name: "Name",
            endpoint: "Pick endpoint",
            configtype: "Pick configuration type",
            create: "Create configuration"
        },
    },
    donate: {
        title: "Support Project Development",
        description: "This project is created and maintained by me in my free time with enthusiasm and remains completely free.\nYour support will help keep the project functioning properly and allow more time for development, improvement, and implementation of new features.",
        buttonYoumoney: "Support",
        buttonTips: "Support via Tips",
    },
    user: {
        quickactions: {
            create: "Create configuration",
            allconf: "All configurations",
            profile: "Profile"
        },
        profile: {
            title: "Information profile",
            user: "User",
            created: "Created by",
            conf: "Configuration",
            password: {
                title: "Change password",
                current: "Current password",
                new: "New password",
                new_hint: "Minimum 3 symbols",
                confirm: "Confirm password",
                change: "Change password"
            }
        }
    }
};