📦 Kopsis-Spensa

Kopsis-Spensa adalah aplikasi berbasis Golang yang dirancang untuk mendukung sistem koperasi sekolah dengan arsitektur modular, aman, dan mudah dikembangkan. Proyek ini menggunakan PostgreSQL sebagai database utama dan Redis sebagai cache/session store, semuanya berjalan dalam container Docker untuk memudahkan deployment dan scaling.

✨ Fitur Utama

1. Backend Golang dengan struktur proyek yang terorganisir (cmd, internal, handlers, models, dll).

2. Database PostgreSQL untuk penyimpanan data yang konsisten dan reliabel.

3. Redis sebagai caching layer untuk meningkatkan performa aplikasi.

4. Docker Compose untuk orkestrasi container (app, database, redis).

5. Konfigurasi Environment (.env) yang fleksibel untuk memudahkan pengaturan.

6. Integrasi dengan aaPanel untuk manajemen server yang lebih mudah (domain, SSL, monitoring).

🛠️ Teknologi yang Digunakan

    Go (Golang) – bahasa pemrograman utama.

    PostgreSQL – relational database.

    Redis – in-memory data store.

    Docker & Docker Compose – containerization & orchestration.

    aaPanel – server management panel.

    Nginx/Reverse Proxy – untuk routing domain dan SSL.