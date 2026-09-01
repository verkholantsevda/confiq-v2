INSERT INTO config_types (
    name,
    description,
    config_template,
    usage_instructions,
    client_links,
    is_active
)
SELECT
    'AmneziaWG',
    'WireGuard с обфускацией трафика. Помогает обходить DPI блокировки.',
    '[Interface]
PrivateKey = {{ private_key }}
Address = {{ client_ipv4 }}/32, {{ client_ipv6 }}/128
DNS = 1.1.1.1, 1.0.0.1, 2606:4700:4700::1111, 2606:4700:4700::1001
MTU = 1280
Jc = 120
Jmin = 23
Jmax = 911
S1 = 0
S2 = 0
H1 = 1
H2 = 2
H3 = 3
H4 = 4

[Peer]
PublicKey = {{ peer_public_key }}
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = {{ endpoint }}:{{ port }}
PersistentKeepalive = 25
',
    '<div class="usage-instructions">
    <h4>Инструкция по использованию AmneziaWG</h4>

    <h5>1. Установите клиент AmneziaWG:</h5>
    <ul>
        <li><strong>Windows:</strong> Скачайте с <a href="https://github.com/amnezia-vpn/amneziawg-windows-client/releases" target="_blank">GitHub Releases</a></li>
        <li><strong>MacOS:</strong> App Store <a href="https://apps.apple.com/ru/app/amneziawg/id6478942365" target="_blank">AmneziaWG</a></li>
        <li><strong>iOS:</strong> App Store <a href="https://apps.apple.com/ru/app/amneziawg/id6478942365" target="_blank">AmneziaWG</a></li>
        <li><strong>Android:</strong> Google Play <a href="https://github.com/amnezia-vpn/amneziawg-android/releases" target="_blank">AmneziaWG</a></li>
        <li><strong>Linux:</strong> Следуйте инструкциям на <a href="https://github.com/amnezia-vpn/amneziawg-linux-kernel-module" target="_blank">GitHub</a></li>
        <li style="list-style: none; margin-top: 12px;">
            <strong>Если AmneziaWG отсутствует в вашем регионе, воспользуйтесь DefaultVPN для Apple:</strong>
        </li>
        <li><strong>macOS:</strong> App Store <a href="https://apps.apple.com/ru/app/defaultvpn/id6744725017" target="_blank">DefaultVPN</a></li>
        <li><strong>iOS:</strong> App Store <a href="https://apps.apple.com/ru/app/defaultvpn/id6744725017" target="_blank">DefaultVPN</a></li>
    </ul>

    <h5>2. Что такое AmneziaWG?</h5>
    <p>AmneziaWG - это форк WireGuard с добавлением обфускации трафика. Параметры Jc, Jmin, Jmax добавляют "шум" в пакеты, что помогает обходить блокировки DPI (Deep Packet Inspection).</p>

    <h5>3. Импорт конфигурации:</h5>
    <ul>
        <li>Скопируйте содержимое конфигурации</li>
        <li>Откройте AmneziaWG клиент</li>
        <li>Нажмите "Import tunnel from file"</li>
        <li>Выберите скачанный файл</li>
        <li>Или добавьте конфигурацию через QR Code</li>
    </ul>

    <h5>4. Подключение:</h5>
    <ul>
        <li>Нажмите "Activate" для подключения</li>
        <li>Значок туннеля станет зелёным при успешном подключении</li>
    </ul>

    <strong>Важно:</strong> AmneziaWG несовместим с обычным WireGuard. Убедитесь, что на всех устройствах используется AmneziaWG клиент.
</div>',
    '{
        "windows": "https://github.com/amnezia-vpn/amneziawg-windows-client/releases",
        "macos": "https://apps.apple.com/ru/app/amneziawg/id6478942365",
        "ios": "https://apps.apple.com/ru/app/amneziawg/id6478942365",
        "android": "https://play.google.com/store/apps/details?id=org.amnezia.awg",
        "linux": "https://github.com/amnezia-vpn/amneziawg-linux-kernel-module"
    }',
    1
WHERE NOT EXISTS (
    SELECT 1 FROM config_types WHERE name = 'AmneziaWG'
);