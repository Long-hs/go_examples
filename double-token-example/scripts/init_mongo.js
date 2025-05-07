// 创建数据库
db = db.getSiblingDB('double_token_example');

// 创建商品集合
db.createCollection('goods', {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["name", "price", "stock", "start_time", "end_time"],
            properties: {
                name: {
                    bsonType: "string",
                    description: "商品名称"
                },
                price: {
                    bsonType: "double",
                    description: "商品价格"
                },
                stock: {
                    bsonType: "int",
                    description: "商品库存"
                },
                start_time: {
                    bsonType: "date",
                    description: "秒杀开始时间"
                },
                end_time: {
                    bsonType: "date",
                    description: "秒杀结束时间"
                },
                created_at: {
                    bsonType: "date",
                    description: "创建时间"
                },
                updated_at: {
                    bsonType: "date",
                    description: "更新时间"
                }
            }
        }
    }
});

// 创建索引
db.goods.createIndex({ "name": 1 });
db.goods.createIndex({ "start_time": 1 });
db.goods.createIndex({ "end_time": 1 });

print("MongoDB 初始化完成！"); 