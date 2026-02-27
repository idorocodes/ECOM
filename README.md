# Ecom: The Zero-Dependency E-commerce Server

A high-performance, lightweight e-commerce backend built entirely with the Go Standard Library. No external drivers, no frameworks, no ORMs—just pure Go.

##  The "Pure Go" Philosophy
Modern Go development is often cluttered with hundreds of node_modules-style dependencies. This project proves that you can build a scalable, secure, and maintainable e-commerce platform using only what the Go team provided in the box.

**Binary Size:** Extremely small and portable.

**Security:** Audit-friendly; every line of code is either yours or the Go core team's.

**Performance:** Minimal memory overhead without the "magic" of heavy frameworks.

##  Core Features
###  Authentication & Security
- Manual JWT: Custom implementation of HMAC-SHA256 signing.
- Security-First: Custom Password Hashing using crypto/sha512 and constant-time comparisons.
- Zero-Dep Middleware: Auth guards for protected routes (Checkout, Profile).

###  Product & Inventory Management
- RESTful API: Full CRUD operations for products and categories.
- Supabase REST Integration: Utilizing the PostgREST interface via net/http to manage persistent data.
- Concurrency: Uses Go's native context and http pooling for high-traffic scalability.

###  Shopping Cart & Orders
- State Management: Efficient handling of cart operations.
- Order Processing: Transactional logic flow built into the service layer.

##  Getting Started
1. **Database Setup**
Ensure your Supabase/Postgres instance has the required tables (users, products, orders). Ensure Row Level Security (RLS) is configured to allow your API key to perform the necessary operations.
2. **Configuration**
The system expects the following environment variables:
  - `SUPA_URL`: Your Supabase Project REST URL.
  - `SUPA_KEY`: Your Anon/Service Role Key.
  - `JWT_SECRET`: The secret key for signing tokens.
3. **Run**
bash 
# No go get needed!
go run cmd/main.go


 
##  License 
This project is open-source and intended for those who value simplicity and the power of the Go Standard Library.
