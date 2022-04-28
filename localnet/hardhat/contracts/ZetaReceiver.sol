// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "./ZetaInterfaces.sol";

interface ZetaReceiver {
    function onZetaMessage(ZetaInterfaces.ZetaMessage calldata zetaMessage)
        external;

    function onZetaRevert(ZetaInterfaces.ZetaRevert calldata zetaRevert)
        external;
}
