// Centralized max length limits for text inputs
// Adjust here to tune lengths consistently across the app

export const LENGTHS = {
  traceCode: 18,
  farmer: {
    fruitName: 100,
    origin: 100,
    farmerName: 100,
    supplierPhone: 30
  },
  factory: {
    productName: 100,
    batch: 32,
    factoryName: 200,
    contactNumber: 30
  },
  driver: {
    name: 50,
    phone: 30,
    carNumber: 10,
    transport: 200
  },
  shop: {
    name: 100,
    address: 200,
    phone: 30
  }
}
