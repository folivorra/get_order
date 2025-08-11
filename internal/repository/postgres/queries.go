package postgres

const (
	existsCheckQuery = `
	SELECT EXISTS (SELECT order_uid FROM orders WHERE order_uid = $1)
	`
	deliverySaveQuery = `
	INSERT INTO deliveries (
		delivery_uid, name, phone, zip, city, address, region, email
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	);
	`
	paymentSaveQuery = `
	INSERT INTO payments (
	    payment_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
	);
	`
	orderSaveQuery = `
	INSERT INTO orders (
		order_uid, track_number, entry, delivery_uid, payment_uid, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	);
	`
	itemSaveQuery = `
	INSERT INTO items (
	    item_uid, chrt_id, track_number, rid, name, size, nm_id, brand, status
	) VALUES (
	    $1, $2, $3, $4, $5, $6, $7, $8, $9
	)
	ON CONFLICT (item_uid) DO NOTHING;
	`
	itemOrderSaveQuery = `
	INSERT INTO order_item (
	    order_item_uid, order_uid, item_uid, price, sale, total_price, quantity
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	);
	`
	orderGetQuery = `
	SELECT *
	FROM orders o
	JOIN deliveries d ON d.delivery_uid = o.delivery_uid
	JOIN payments p ON p.payment_uid = o.payment_uid
	JOIN order_item oi ON oi.order_uid = o.order_uid
	JOIN items i ON oi.item_uid = i.item_uid
	WHERE o.order_uid = $1;
	`
)
