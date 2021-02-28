// const mongo = require('mongodb')
// const axios = require('axios')
// const chanceInstance = new chance()
// const Config = require("./config")

// async function sendToSearch(sProducts) {
//     await axios
//     .post(Config.ProductBulkInsertURL, sProducts)
//     .then(res => {
//       console.log(`statusCode: ${res.statusCode}`)
//     })
//     .catch(error => {
//       console.error(error.response.data)
//     })
// }

// function CtoSProduct(counter) {
//     return {
//         id : counter,
//         shop_item_id: cProduct.shop_item_id,
//         version: cProduct.version,
//         slug: cProduct.shop_item_slug,
//         name: cProduct.item_name,
//         shop_name: cProduct.shop_name,
//         shop_slug: cProduct.shop_slug,
//         brand_name: cProduct.brand_name,
//         brand_slug: cProduct.brand_slug,
//         category_name: cProduct.category_name,
//         category_slug: cProduct.category_slug,
//         color: cProduct.color,
//         color_variants: cProduct.color_variants,
//         discounted_price: cProduct.discounted_price,
//         max_price: cProduct.max_price,
//         min_price: cProduct.min_price,
//         price: cProduct.item_price,
//         product_image: cProduct.item_images ? cProduct.item_images.length ? cProduct.item_images[0]: "" : "",
//     }
// }

// (async function main() {
//     let counter = 1
//     let data = []
//     while (counter < 50000) {
//         data = [...data, CtoSProduct(counter)]
//         console.log(counter)
//         if (counter % 500 == 0) {
//             await sendToSearch(data)
//             console.log(counter, " products inserted")
//             data = []
//         }
       
//         counter += 1
//     }
//     console.log("finished bulk inserting products")
//     process.exit(0)
// })()

