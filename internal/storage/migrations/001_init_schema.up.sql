-- Сущности
CREATE TABLE IF NOT EXISTS sto (
        id TEXT PRIMARY KEY,
        name TEXT,
        ext_index TEXT,
        address TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

CREATE TABLE IF NOT EXISTS warehouse (
        id TEXT PRIMARY KEY,
        name TEXT,
        ext_index TEXT,
        address TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

CREATE TABLE IF NOT EXISTS client (
        id TEXT PRIMARY KEY,
        type TEXT,
        short_name TEXT,
        unp TEXT,
        full_name TEXT,
        address TEXT,
        bank_name TEXT,
        bik TEXT,
        bank_account TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

CREATE TABLE IF NOT EXISTS contract (
        id TEXT PRIMARY KEY,
        client_id TEXT,
        number TEXT,
        date TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

CREATE TABLE IF NOT EXISTS car (
        id TEXT PRIMARY KEY,
        client_id TEXT,
        brand TEXT,
        model TEXT,
        color TEXT,
        year TEXT,
        license_plate_number TEXT,
        engine_code TEXT,
        vin TEXT,
        mileage TEXT,
        comment TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

CREATE TABLE IF NOT EXISTS product (
        id TEXT PRIMARY KEY,
        art_id TEXT,
        pin TEXT,
        name TEXT,
        description TEXT,
        unit TEXT,
        percent_vat TEXT,
        ttn_date TEXT,
        ttn_number TEXT,
        cost TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );



-- Документы
CREATE TABLE IF NOT EXISTS invoice (
        id TEXT PRIMARY KEY,
        ext_index TEXT,
        doc_date TEXT,
        doc_number TEXT,
        status BOOLEAN,
        sto_id TEXT,
        warehouse_id TEXT,
        client_id TEXT,
        sum TEXT,
        sum_vat TEXT,
        comment TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

    CREATE TABLE IF NOT EXISTS invoice_position (
            id TEXT,
            invoice_id TEXT,
            product_id TEXT,
            amount TEXT,
            price TEXT,
            price_without_vat TEXT,
            sum_without_vat TEXT,
            percent_vat TEXT,
            sum_vat TEXT,
            sum TEXT,
            PRIMARY KEY (id, invoice_id)
        );


CREATE TABLE IF NOT EXISTS realization (
        id TEXT PRIMARY KEY,
        ext_index TEXT,
        doc_date TEXT,
        doc_number TEXT,
        status BOOLEAN,
        sto_id TEXT,
        client_id TEXT,
        contract_id TEXT,
        sum TEXT,
        sum_vat TEXT,
        comment TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

    CREATE TABLE IF NOT EXISTS realization_position (
            id TEXT,
            realization_id TEXT,
            product_id TEXT,
            amount TEXT,
            price TEXT,
            price_without_vat TEXT,
            sum_without_vat TEXT,
            percent_vat TEXT,
            sum_vat TEXT,
            sum TEXT,
            PRIMARY KEY (id, realization_id)
        );


CREATE TABLE IF NOT EXISTS moving (
        id TEXT PRIMARY KEY,
        ext_index TEXT,
        doc_date TEXT,
        doc_number TEXT,
        status BOOLEAN,
        sto_from_id TEXT,
        sto_to_id TEXT,
        warehouse_from_id TEXT,
        warehouse_to_id TEXT,
        sum TEXT,
        comment TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

    CREATE TABLE IF NOT EXISTS moving_position (
            id TEXT,
            moving_id TEXT,
            product_id TEXT,
            amount TEXT,
            price TEXT,
            sum TEXT,
            PRIMARY KEY (id, moving_id)
        );


CREATE TABLE IF NOT EXISTS dismantling (
        id TEXT PRIMARY KEY,
        ext_index TEXT,
        doc_date TEXT,
        doc_number TEXT,
        status BOOLEAN,
        sto_id TEXT,
        warehouse_id TEXT,
        product_id TEXT,
        amount TEXT,
        sum TEXT,
        comment TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

    CREATE TABLE IF NOT EXISTS dismantling_position (
            id TEXT,
            dismantling_id TEXT,
            product_id TEXT,
            amount TEXT,
            price TEXT,
            sum TEXT,
            PRIMARY KEY (id, dismantling_id)
        );


CREATE TABLE IF NOT EXISTS inventory (
        id TEXT PRIMARY KEY,
        ext_index TEXT,
        doc_date TEXT,
        doc_number TEXT,
        status BOOLEAN,
        sto_id TEXT,
        warehouse_id TEXT,
        sum TEXT,
        comment TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

    CREATE TABLE IF NOT EXISTS inventory_position (
            id TEXT,
            inventory_id TEXT,
            product_id TEXT,
            amount TEXT,
            amount_plan TEXT,
            price TEXT,
            sum TEXT,
            PRIMARY KEY (id, inventory_id)
        );


CREATE TABLE IF NOT EXISTS request (
        id TEXT PRIMARY KEY,
        ext_index TEXT,
        doc_date TEXT,
        doc_number TEXT,
        status BOOLEAN,
        sto_id TEXT,
        warehouse_id TEXT,
        client_id TEXT,
        contract_id TEXT,
        car_id TEXT,
        garant_work TEXT,
        legal_person_short_name TEXT,
        client_individual_name TEXT,
        client_legal_short_name TEXT,
        contact_person_fio TEXT,
        car_brand TEXT,
        car_engine_capacity TEXT,
        car_license_plate_number TEXT,
        car_manufacture_year TEXT,
        car_mileage TEXT,
        car_model_name TEXT,
        car_modification TEXT,
        car_vin TEXT,
        complete_status_date TEXT,
        close_status_date TEXT,
        sum_work TEXT,
        sum_parts TEXT,
        sum_req TEXT,
        created_user_short_fio TEXT,
        normo_time_fact TEXT,
        normo_time_plan TEXT,
        percent_vat TEXT,
        sum_req_without_vat TEXT,
        sum_vat TEXT,
        reason_for_petition TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_updated BOOLEAN,
        is_fixed BOOLEAN DEFAULT FALSE,
        base_uid TEXT DEFAULT ""
    );

    CREATE TABLE IF NOT EXISTS request_position (
            id TEXT,
            request_id TEXT,
            product_id TEXT,
            amount TEXT,
            price TEXT,
            price_without_vat TEXT,
            sum_without_vat TEXT,
            percent_vat TEXT,
            sum_vat TEXT,
            sum TEXT,
            client_parts BOOLEAN,
            used_parts BOOLEAN,
            PRIMARY KEY (id, request_id)
        );

    CREATE TABLE IF NOT EXISTS request_work (
            id TEXT,
            request_id TEXT,
            work_name TEXT,
            amount TEXT,
            price TEXT,
            price_without_vat TEXT,
            sum_without_vat TEXT,
            percent_vat TEXT,
            sum_vat TEXT,
            sum TEXT,
            PRIMARY KEY (id, request_id)
        );

    CREATE TABLE IF NOT EXISTS request_performer (
            id TEXT,
            request_id TEXT,
            work_id TEXT,
            name TEXT,
            cost_per_hour TEXT,
            total_earnings TEXT,
            productivity_percent TEXT,
            PRIMARY KEY (id, request_id)
        );



-- Таблица хэшей
CREATE TABLE IF NOT EXISTS entity_hashes (
        object_type TEXT NOT NULL,
        entity_id TEXT NOT NULL,
        hash TEXT NOT NULL,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (object_type, entity_id)
    );

