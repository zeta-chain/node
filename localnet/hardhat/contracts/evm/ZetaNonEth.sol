/* solhint-disable var-name-mixedcase */
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";

contract ZetaNonEth is ERC20Burnable {
    /**
     * @dev Collectively hold by Zeta blockchain validators
     */
    address public TSSAddress;

    /**
     * @dev Initially a multi-sig, eventually hold by Zeta blockchain validators (via renounceTSSAddressUpdater)
     */
    address public TSSAddressUpdater;

    address public MPIAddress;

    event MMinted(
        address indexed mintee,
        uint amount,
        bytes32 indexed internalSendHash
    );
    event MBurnt(address indexed burnee, uint amount);

    constructor(
        uint256 initialSupply,
        address _TSSAddress,
        address _TSSAddressUpdater
    ) ERC20("Zeta", "ZETA") {
        _mint(msg.sender, initialSupply * (10**uint256(decimals())));

        TSSAddress = _TSSAddress;
        TSSAddressUpdater = _TSSAddressUpdater;
    }

    function updateTSSAndMPIAddresses(address _tss, address _mpi) external {
        require(
            msg.sender == TSSAddressUpdater,
            "ZetaNonEth: sender is not TSSAddressUpdater"
        );

        TSSAddress = _tss;
        MPIAddress = _mpi;
    }

    /**
     * @dev Sets TSSAddressUpdater to be TSSAddress
     */
    function renounceTSSAddressUpdater() external {
        require(
            msg.sender == TSSAddressUpdater,
            "ZetaNonEth: sender is not TSSAddressUpdater"
        );
        require(TSSAddress != address(0), "ZetaNonEth: Invalid TSSAddress");

        TSSAddressUpdater = TSSAddress;
    }

    function mint(
        address mintee,
        uint value,
        bytes32 internalSendHash
    ) external {
        /**
         * @dev Only MPI or TSS can mint since minting requires burning the equivalent in another chain
         */
        require(
            msg.sender == TSSAddress || msg.sender == MPIAddress,
            "ZetaNonEth: only TSSAddress or MPIAddress can mint"
        );

        _mint(mintee, value);

        emit MMinted(mintee, value, internalSendHash);
    }
}
/* solhint-enable var-name-mixedcase */
