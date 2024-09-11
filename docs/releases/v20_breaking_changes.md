
# V20 Breaking Changes

### Emissions factors deprecated

* `EmissionsFactors` have been deprecated and removed from the `emissions` module. 
  - This results in the removal of the query `/zeta-chain/emissions/get_emissions_factors`.
  - The fixed block reward amount can now be queried via `/zeta-chain/emissions/params`. This is constant for every block and does not depend on any factors.

