import type { Endpoint } from './endpoint';
import type { ConfigType } from './configType';

export interface Configuration {
    id: number;
    name: string;
    user_id: number;
    endpoint_id?: number;
    config_type_id?: number;
    config_type: {
        id: number;
        name: string;
        usage_instructions: string;
    };
    endpoint: {
        id: number;
        name: string;
        address: string;
        port: number;
    };
    config_content: string;
    public_key: string;
    peer_public_key: string;
    client_ipv4: string;
    client_ipv6: string;
    created_at: string;
}