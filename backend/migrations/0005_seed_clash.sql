INSERT INTO config_types (
    name,
    description,
    config_template,
    usage_instructions,
    client_links,
    is_active
)
SELECT
    'Clash',
    'Clash/Mihomo прокси-клиент с поддержкой правил маршрутизации и WireGuard.',
    'mixed-port: 7890
allow-lan: false
bind-address: ''*''
mode: rule
log-level: info
external-controller: 127.0.0.1:9090
unified-delay: true

profile:
  store-selected: true
  store-fake-ip: true

dns:
  enable: true
  listen: 0.0.0.0:1053
  default-nameserver:
    - 223.5.5.5
    - 8.8.8.8
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  use-hosts: true
  nameserver:
    - https://doh.pub/dns-query
    - https://dns.alidns.com/dns-query
  fallback:
    - https://1.1.1.1/dns-query
    - https://dns.cloudflare.com/dns-query
  fallback-filter:
    geoip: true
    geoip-code: CN
    ipcidr:
      - 240.0.0.0/4

proxies:
  - name: "WARP-{{ endpoint }}"
    type: wireguard
    server: {{ endpoint }}
    port: {{ port }}
    ip: {{ client_ipv4 }}
    ipv6: {{ client_ipv6 }}
    private-key: {{ private_key }}
    public-key: {{ peer_public_key }}
    udp: true
    reserved: []

proxy-groups:
  - name: "PROXY"
    type: select
    proxies:
      - "WARP-{{ endpoint }}"
      - "DIRECT"

rules:
  - GEOIP,CN,DIRECT
  - GEOSITE,cn,DIRECT
  - MATCH,PROXY
',
    '<div class="usage-instructions">
    <h4>Инструкция по использованию Clash</h4>

    <h5>1. Установите клиент Clash:</h5>
    <ul>
        <li><strong>Windows:</strong> Clash Verge Rev, Clash for Windows, или Clash Verge</li>
        <li><strong>macOS:</strong> Clash Verge Rev, ClashX, ClashX Pro</li>
        <li><strong>iOS:</strong> Stash, Shadowrocket (платные)</li>
        <li><strong>Android:</strong> Clash for Android, Clash Meta for Android</li>
        <li><strong>Linux:</strong> Clash Verge Rev</li>
    </ul>

    <h5>2. Рекомендуемые клиенты:</h5>
    <table class="table table-sm">
        <tr>
            <td><strong>Windows/macOS/Linux:</strong></td>
            <td><a href="https://github.com/clash-verge-rev/clash-verge-rev/releases" target="_blank">Clash Verge Rev</a> (бесплатный, рекомендуется)</td>
        </tr>
        <tr>
            <td><strong>iOS:</strong></td>
            <td>Stash (App Store, ~$3.99)</td>
        </tr>
        <tr>
            <td><strong>Android:</strong></td>
            <td><a href="https://github.com/MetaCubeX/ClashMetaForAndroid/releases" target="_blank">Clash Meta for Android</a></td>
        </tr>
    </table>

    <h5>3. Импорт конфигурации:</h5>
    <ul>
        <li>Скачайте файл конфигурации (.yaml или .yml)</li>
        <li>В Clash клиенте нажмите "Profiles" → "Import"</li>
        <li>Выберите скачанный файл или перетащите его в окно</li>
        <li>Активируйте профиль кликом по нему</li>
    </ul>

    <h5>4. Настройка системного прокси:</h5>
    <ul>
        <li>Включите "System Proxy" в главном окне Clash</li>
        <li>Или настройте приложения на использование HTTP прокси 127.0.0.1:7890</li>
    </ul>

    <h5>5. Управление через веб-интерфейс:</h5>
    <ul>
        <li>Откройте http://127.0.0.1:9090/ui в браузере</li>
        <li>Или используйте встроенный интерфейс клиента</li>
    </ul>

    <strong>Преимущества Clash:</strong> Поддержка правил маршрутизации (rules), автоматическое переключение прокси, встроенный DNS с fake-ip.
</div>',
    '{
        "windows": "https://github.com/clash-verge-rev/clash-verge-rev/releases",
        "macos": "https://github.com/clash-verge-rev/clash-verv-rev/releases",
        "ios": "https://apps.apple.com/us/app/stash/id1582679995",
        "android": "https://github.com/MetaCubeX/ClashMetaForAndroid/releases",
        "linux": "https://github.com/clash-verge-rev/clash-verge-rev/releases"
    }',
    1
WHERE NOT EXISTS (
    SELECT 1 FROM config_types WHERE name = 'Clash'
);