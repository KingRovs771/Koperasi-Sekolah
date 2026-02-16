--
-- PostgreSQL database dump
--

\restrict EUy8ab1yZJ8nIbp5c4quabtaHpIHizH5r5l21r3Ln17BZHD1BDsUibhTDzSUfJX

-- Dumped from database version 18.1 (Debian 18.1-1.pgdg13+2)
-- Dumped by pg_dump version 18.1 (Debian 18.1-1.pgdg13+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP INDEX public.idx_users_users_id;
DROP INDEX public.idx_users_deleted_at;
DROP INDEX public.idx_transactions_transactions_id;
DROP INDEX public.idx_transactions_items_transactions_items_id;
DROP INDEX public.idx_products_products_uid;
DROP INDEX public.idx_products_products_id;
DROP INDEX public.idx_products_category_uid;
DROP INDEX public.idx_mutasi_stocks_mutasi_stocks_id;
DROP INDEX public.idx_categories_category_id;
ALTER TABLE ONLY public.users DROP CONSTRAINT users_pkey;
ALTER TABLE ONLY public.transactions DROP CONSTRAINT uni_transactions_user_uid;
ALTER TABLE ONLY public.transactions DROP CONSTRAINT uni_transactions_transactions_uid;
ALTER TABLE ONLY public.transactions DROP CONSTRAINT uni_transactions_no_faktur;
ALTER TABLE ONLY public.transactions_items DROP CONSTRAINT uni_transactions_items_transactions_uid;
ALTER TABLE ONLY public.products DROP CONSTRAINT uni_products_kode_produk;
ALTER TABLE ONLY public.categories DROP CONSTRAINT uni_categories_category_uid;
ALTER TABLE ONLY public.transactions DROP CONSTRAINT transactions_pkey;
ALTER TABLE ONLY public.transactions_items DROP CONSTRAINT transactions_items_pkey;
ALTER TABLE ONLY public.products DROP CONSTRAINT products_pkey;
ALTER TABLE ONLY public.mutasi_stocks DROP CONSTRAINT mutasi_stocks_pkey;
ALTER TABLE ONLY public.categories DROP CONSTRAINT categories_pkey;
ALTER TABLE ONLY public.backup_restores DROP CONSTRAINT backup_restores_pkey;
ALTER TABLE public.users ALTER COLUMN users_id DROP DEFAULT;
ALTER TABLE public.transactions_items ALTER COLUMN transactions_items_id DROP DEFAULT;
ALTER TABLE public.transactions ALTER COLUMN transactions_id DROP DEFAULT;
ALTER TABLE public.products ALTER COLUMN products_id DROP DEFAULT;
ALTER TABLE public.mutasi_stocks ALTER COLUMN mutasi_stocks_id DROP DEFAULT;
ALTER TABLE public.categories ALTER COLUMN category_id DROP DEFAULT;
ALTER TABLE public.backup_restores ALTER COLUMN backup_restore_id DROP DEFAULT;
DROP SEQUENCE public.users_users_id_seq;
DROP TABLE public.users;
DROP SEQUENCE public.transactions_transactions_id_seq;
DROP SEQUENCE public.transactions_items_transactions_items_id_seq;
DROP TABLE public.transactions_items;
DROP TABLE public.transactions;
DROP SEQUENCE public.products_products_id_seq;
DROP TABLE public.products;
DROP SEQUENCE public.mutasi_stocks_mutasi_stocks_id_seq;
DROP TABLE public.mutasi_stocks;
DROP SEQUENCE public.categories_category_id_seq;
DROP TABLE public.categories;
DROP SEQUENCE public.backup_restores_backup_restore_id_seq;
DROP TABLE public.backup_restores;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: backup_restores; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.backup_restores (
    backup_restore_id bigint NOT NULL,
    backup_restore_uid text,
    nama_file text,
    ukuran_file bigint,
    created_at timestamp with time zone
);


ALTER TABLE public.backup_restores OWNER TO postgres;

--
-- Name: backup_restores_backup_restore_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.backup_restores_backup_restore_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.backup_restores_backup_restore_id_seq OWNER TO postgres;

--
-- Name: backup_restores_backup_restore_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.backup_restores_backup_restore_id_seq OWNED BY public.backup_restores.backup_restore_id;


--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    category_id bigint NOT NULL,
    category_uid text,
    nama_category text NOT NULL,
    icon character varying(50) DEFAULT 'box'::character varying,
    description text,
    status text DEFAULT 'active'::text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: categories_category_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_category_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_category_id_seq OWNER TO postgres;

--
-- Name: categories_category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_category_id_seq OWNED BY public.categories.category_id;


--
-- Name: mutasi_stocks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mutasi_stocks (
    mutasi_stocks_id bigint NOT NULL,
    mutasi_stocks_uid text,
    products_uid text,
    user_uid text,
    tipe text NOT NULL,
    jumlah bigint,
    catatan text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.mutasi_stocks OWNER TO postgres;

--
-- Name: mutasi_stocks_mutasi_stocks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.mutasi_stocks_mutasi_stocks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.mutasi_stocks_mutasi_stocks_id_seq OWNER TO postgres;

--
-- Name: mutasi_stocks_mutasi_stocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.mutasi_stocks_mutasi_stocks_id_seq OWNED BY public.mutasi_stocks.mutasi_stocks_id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    products_id bigint NOT NULL,
    products_uid text,
    kode_produk character varying(50),
    category_uid text,
    nama_products text,
    description text,
    harga_jual bigint,
    harga_beli bigint,
    stock bigint,
    stock_min bigint,
    image text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.products OWNER TO postgres;

--
-- Name: products_products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_products_id_seq OWNER TO postgres;

--
-- Name: products_products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_products_id_seq OWNED BY public.products.products_id;


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions (
    transactions_id bigint NOT NULL,
    transactions_uid text,
    no_faktur text NOT NULL,
    user_uid text NOT NULL,
    total_belanja numeric NOT NULL,
    total_modal numeric NOT NULL,
    metode_pembayaran text DEFAULT 'tunai'::text,
    uang_diterima numeric,
    kembalian numeric,
    created_at timestamp with time zone
);


ALTER TABLE public.transactions OWNER TO postgres;

--
-- Name: transactions_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions_items (
    transactions_items_id bigint NOT NULL,
    transactions_uid text,
    products_uid text NOT NULL,
    jumlah bigint NOT NULL,
    harga_saat_ini numeric NOT NULL,
    modal_saat_ini numeric NOT NULL,
    sub_total numeric NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.transactions_items OWNER TO postgres;

--
-- Name: transactions_items_transactions_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.transactions_items_transactions_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.transactions_items_transactions_items_id_seq OWNER TO postgres;

--
-- Name: transactions_items_transactions_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.transactions_items_transactions_items_id_seq OWNED BY public.transactions_items.transactions_items_id;


--
-- Name: transactions_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.transactions_transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.transactions_transactions_id_seq OWNER TO postgres;

--
-- Name: transactions_transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.transactions_transactions_id_seq OWNED BY public.transactions.transactions_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    users_id bigint NOT NULL,
    users_uid text,
    nama_lengkap text NOT NULL,
    username text,
    email text,
    password text,
    no_hp text,
    alamat text,
    role text DEFAULT 'kasir'::text,
    status text DEFAULT 'active'::text,
    foto text,
    created_at timestamp with time zone,
    update_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_users_id_seq OWNER TO postgres;

--
-- Name: users_users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_users_id_seq OWNED BY public.users.users_id;


--
-- Name: backup_restores backup_restore_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.backup_restores ALTER COLUMN backup_restore_id SET DEFAULT nextval('public.backup_restores_backup_restore_id_seq'::regclass);


--
-- Name: categories category_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN category_id SET DEFAULT nextval('public.categories_category_id_seq'::regclass);


--
-- Name: mutasi_stocks mutasi_stocks_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mutasi_stocks ALTER COLUMN mutasi_stocks_id SET DEFAULT nextval('public.mutasi_stocks_mutasi_stocks_id_seq'::regclass);


--
-- Name: products products_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN products_id SET DEFAULT nextval('public.products_products_id_seq'::regclass);


--
-- Name: transactions transactions_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions ALTER COLUMN transactions_id SET DEFAULT nextval('public.transactions_transactions_id_seq'::regclass);


--
-- Name: transactions_items transactions_items_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions_items ALTER COLUMN transactions_items_id SET DEFAULT nextval('public.transactions_items_transactions_items_id_seq'::regclass);


--
-- Name: users users_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN users_id SET DEFAULT nextval('public.users_users_id_seq'::regclass);


--
-- Data for Name: backup_restores; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.backup_restores (backup_restore_id, backup_restore_uid, nama_file, ukuran_file, created_at) FROM stdin;
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (category_id, category_uid, nama_category, icon, description, status, created_at, updated_at) FROM stdin;
1	cc2157e6-6598-4f01-9d9b-5f2204f6d24e	Alat Tulis	Pencil	Berisi Alat Tulis Seperti Pensil, Penggaris, Pulpen, Penghapus dan lain sebagainya	active	2026-02-03 05:07:01.116234+00	2026-02-03 05:07:01.116234+00
2	8a4ed5c0-7cc8-4678-9c7b-e4fa61abe2d9	Seragam	Cloth	Berisi Perlengkapan Seragam seperti topi, dasi, ikat pinggang dan lain sebagainya, 	active	2026-02-03 05:07:01.139754+00	2026-02-03 05:07:01.139754+00
3	bb8bf073-e0f9-451f-ae0d-464f2008ba45	Pramuka	Scout	Berisi Alat Pramuka seperti semaphore, Peluit, Tali, dan lain sebagainya 	active	2026-02-03 05:07:01.14466+00	2026-02-03 05:07:01.14466+00
\.


--
-- Data for Name: mutasi_stocks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.mutasi_stocks (mutasi_stocks_id, mutasi_stocks_uid, products_uid, user_uid, tipe, jumlah, catatan, created_at) FROM stdin;
1	7b00483d-1b68-481d-a706-ca75294ca5e8	4ffc439b-f7b2-4414-901e-b4edc8b63a1b	bb59eebf-cbfa-4c5e-b6d3-116225212c60	masuk	4	Restock Bu Arni	2026-02-09 12:35:10.917386+00
2	80011892-59b6-4dad-a2a3-8e6c378b66cb	4ffc439b-f7b2-4414-901e-b4edc8b63a1b	bb59eebf-cbfa-4c5e-b6d3-116225212c60	keluar	27	Barang Sudah tidak layak pakai	2026-02-11 06:58:40.70468+00
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (products_id, products_uid, kode_produk, category_uid, nama_products, description, harga_jual, harga_beli, stock, stock_min, image, created_at, updated_at, deleted_at) FROM stdin;
1	4ffc439b-f7b2-4414-901e-b4edc8b63a1b	BRT01	bb8bf073-e0f9-451f-ae0d-464f2008ba45	Baret Pramuka	Baret Pramuka Cowok	10000	5000	2	4	1770521640_brd-19110_aksesoren-topi-baret-pramuka-rajut-dan-laken-perlengkapan-atribut-pramuka-sekolah-aksesoris-pramuka-laki-laki-topi-cowok-pria-cokelat_full01-c92df39e.jpg	2026-02-08 03:34:00.9354+00	2026-02-11 06:58:40.701285+00	\N
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (transactions_id, transactions_uid, no_faktur, user_uid, total_belanja, total_modal, metode_pembayaran, uang_diterima, kembalian, created_at) FROM stdin;
\.


--
-- Data for Name: transactions_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions_items (transactions_items_id, transactions_uid, products_uid, jumlah, harga_saat_ini, modal_saat_ini, sub_total, created_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (users_id, users_uid, nama_lengkap, username, email, password, no_hp, alamat, role, status, foto, created_at, update_at, deleted_at) FROM stdin;
1	bb59eebf-cbfa-4c5e-b6d3-116225212c60	Rizky Budiarto	kingrovs	kingrovs@smpn1sragen.sch.id	$2a$10$QQ7bJzGpYSLG7QtIPcK41ufmmw9ojIyXyqCB3bFpvuAyjcQGMQXrq	082142896072	Jalan Raya Sukowati No. 162, Sragen Kulon, Sragen	administrator	active		2026-02-09 02:10:35.304257+00	0001-01-01 00:00:00+00	\N
2	642109e9-5cc1-4663-b572-084ca6a6a895	Mukhlis Royyani	siroyy	siroyy@smpn1sragen.sch.id	$2a$10$bGgPeapCZSJJtb53XfNPUOCQAiabgyiq/5c4JxRJdulLCdRKqMFIy	08893884991	Jalan Raya Sukowati No. 162, Sragen Kulon, Sragen	supervisor	active		2026-02-09 02:10:35.355849+00	0001-01-01 00:00:00+00	\N
3	cbf25bc7-ae00-4c20-9545-df538e11b602	Riyana Lili Lestari	riyanalili	riyanalili@smpn1sragen.sch.id	$2a$10$Ya/yo4mwH0szkUEJ/egcSeaCgNbhjnddyZGOVYpnbLoPbFV6h4ZI6	089577389893	Jalan Raya Sukowati No. 162, Sragen Kulon, Sragen	supervisor	active		2026-02-09 02:10:35.404503+00	0001-01-01 00:00:00+00	\N
4	454c139e-822b-416e-8bd5-d730e9894a28	AdminSuper	adminsuperspensa	adminsuperspensa@smpn1sragen.sch.id	$2a$10$U5UA2q8gdTX/XD.u9JZU8urPKb8kpIUQ9cdsTqHB37YbUWxXGc.IS	08779947729	Jalan Raya Sukowati No. 162, Sragen Kulon, Sragen	administrator	active		2026-02-09 02:10:35.455553+00	0001-01-01 00:00:00+00	\N
5	49901716-f2d6-4b79-9a13-609784d4093b	cashier	cashier	cashier@smpn1sragen.sch.id	$2a$10$XF6yhyHM9zpFOhlXhqSRJ.hTSRInX124u6VWdEuAMHYsA6q/LlMIu	08779947729	Jalan Raya Sukowati No. 162, Sragen Kulon, Sragen	kasir	active		2026-02-09 02:10:35.516198+00	0001-01-01 00:00:00+00	\N
6	67be6b35-b62e-43fa-87d1-317fb940cc09	cashier 2	kasirv2	kasirv2@smpn1sragen.sch.id	$2a$10$g0d2a1MdY4jpqtsBCQ9ybuQVApCWTI86LetQ7bP/TsO9XxiZH4wpa			kasir	active		2026-02-09 02:58:19.678575+00	0001-01-01 00:00:00+00	\N
\.


--
-- Name: backup_restores_backup_restore_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.backup_restores_backup_restore_id_seq', 1, false);


--
-- Name: categories_category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_category_id_seq', 3, true);


--
-- Name: mutasi_stocks_mutasi_stocks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.mutasi_stocks_mutasi_stocks_id_seq', 2, true);


--
-- Name: products_products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_products_id_seq', 1, true);


--
-- Name: transactions_items_transactions_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.transactions_items_transactions_items_id_seq', 1, false);


--
-- Name: transactions_transactions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.transactions_transactions_id_seq', 1, false);


--
-- Name: users_users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_users_id_seq', 6, true);


--
-- Name: backup_restores backup_restores_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.backup_restores
    ADD CONSTRAINT backup_restores_pkey PRIMARY KEY (backup_restore_id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (category_id);


--
-- Name: mutasi_stocks mutasi_stocks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mutasi_stocks
    ADD CONSTRAINT mutasi_stocks_pkey PRIMARY KEY (mutasi_stocks_id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (products_id);


--
-- Name: transactions_items transactions_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions_items
    ADD CONSTRAINT transactions_items_pkey PRIMARY KEY (transactions_items_id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (transactions_id);


--
-- Name: categories uni_categories_category_uid; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT uni_categories_category_uid UNIQUE (category_uid);


--
-- Name: products uni_products_kode_produk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT uni_products_kode_produk UNIQUE (kode_produk);


--
-- Name: transactions_items uni_transactions_items_transactions_uid; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions_items
    ADD CONSTRAINT uni_transactions_items_transactions_uid UNIQUE (transactions_uid);


--
-- Name: transactions uni_transactions_no_faktur; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT uni_transactions_no_faktur UNIQUE (no_faktur);


--
-- Name: transactions uni_transactions_transactions_uid; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT uni_transactions_transactions_uid UNIQUE (transactions_uid);


--
-- Name: transactions uni_transactions_user_uid; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT uni_transactions_user_uid UNIQUE (user_uid);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (users_id);


--
-- Name: idx_categories_category_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_categories_category_id ON public.categories USING btree (category_id);


--
-- Name: idx_mutasi_stocks_mutasi_stocks_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_mutasi_stocks_mutasi_stocks_id ON public.mutasi_stocks USING btree (mutasi_stocks_id);


--
-- Name: idx_products_category_uid; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_products_category_uid ON public.products USING btree (category_uid);


--
-- Name: idx_products_products_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_products_products_id ON public.products USING btree (products_id);


--
-- Name: idx_products_products_uid; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_products_products_uid ON public.products USING btree (products_uid);


--
-- Name: idx_transactions_items_transactions_items_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_transactions_items_transactions_items_id ON public.transactions_items USING btree (transactions_items_id);


--
-- Name: idx_transactions_transactions_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_transactions_transactions_id ON public.transactions USING btree (transactions_id);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_users_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_users_users_id ON public.users USING btree (users_id);


--
-- PostgreSQL database dump complete
--

\unrestrict EUy8ab1yZJ8nIbp5c4quabtaHpIHizH5r5l21r3Ln17BZHD1BDsUibhTDzSUfJX

