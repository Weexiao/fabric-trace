# E2E Smoke Test Plan (Manual or Cypress)

This is a minimal smoke test plan to verify main flows:

1. Login Flow
   - Open /login, enter username and password, click Login.
   - Expect dashboard appears and no error toast.

2. Uplink (Manufacturer)
   - Navigate to /uplink as manufacturer.
   - Fill traceability_code with 18 digits, fill product, batch, time, factory name, phone.
   - Choose an image and click Submit.
   - Expect success toast and success dialog appears; copy code works; View Details navigates to /trace/{code}.

3. Role Switch
   - In /uplink, change user role to Dealer.
   - Expect unrelated fields cleared and validation reset.

4. Trace Query
   - Open /trace, input a valid code, click Query.
   - Expect details table shows 4 sections and images load if exist.

If adopting Cypress later, convert into tests under cypress/integration with data-test attributes.

