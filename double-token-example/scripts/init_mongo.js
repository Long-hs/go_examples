// 创建商品集合
db.createCollection("goods", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["name", "price", "stock", "status"],
            properties: {
                name: {
                    bsonType: "string",
                    description: "商品名称"
                },
                description: {
                    bsonType: "string",
                    description: "商品描述"
                },
                price: {
                    bsonType: "double",
                    description: "商品价格"
                },
                stock: {
                    bsonType: "long",
                    description: "商品库存"
                },
                image: {
                    bsonType: "string",
                    description: "商品图片"
                },
                status: {
                    bsonType: "int",
                    description: "商品状态：1-上架，0-下架"
                },
                version: {
                    bsonType: "long",
                    description: "版本号（用于乐观锁）"
                },
                create_time: {
                    bsonType: "date",
                    description: "创建时间"
                },
                update_time: {
                    bsonType: "date",
                    description: "更新时间"
                }
            }
        }
    }
});

// 创建索引
db.goods.createIndex({ "name": 1 });
db.goods.createIndex({ "status": 1 }); 