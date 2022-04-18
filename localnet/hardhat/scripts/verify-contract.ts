import { isZetaAddress, verifyContract } from "../lib/networks";

async function main() {
  const { ZETA_ADDRESS_NAME } = process.env;
  if (!isZetaAddress(ZETA_ADDRESS_NAME))
    throw new Error(`Invalid ZETA_ADDRESS_NAME: ${ZETA_ADDRESS_NAME}`);

  verifyContract(ZETA_ADDRESS_NAME);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
