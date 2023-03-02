
#Grant executer authorization to execute the GAS PRICE VOTER, NONCE VOTER for validator `zeta`
zetacored tx authz grant zeta19wzjdtah4kl2vh77jks68cyy5gpjyurqltys99 generic --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --msg-type=/zetachain.zetacore.crosschain.MsgGasPriceVoter
zetacored tx authz grant zeta19wzjdtah4kl2vh77jks68cyy5gpjyurqltys99 generic --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --msg-type=/zetachain.zetacore.crosschain.MsgNonceVoter

#Grant executer authorization to execute the Inbound and Outbound VOTER for observers `zeta` and `mario`
zetacored tx authz grant zeta19wzjdtah4kl2vh77jks68cyy5gpjyurqltys99 generic --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --msg-type=/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx
zetacored tx authz grant zeta19wzjdtah4kl2vh77jks68cyy5gpjyurqltys99 generic --from mario --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --msg-type=/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx

zetacored tx authz grant zeta19wzjdtah4kl2vh77jks68cyy5gpjyurqltys99 generic --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --msg-type=/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx
zetacored tx authz grant zeta19wzjdtah4kl2vh77jks68cyy5gpjyurqltys99 generic --from mario --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --msg-type=/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx


# Check all the grants
zetacored q authz grants zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk zeta19wzjdtah4kl2vh77jks68cyy5gpjyurqltys99
zetacored q authz grants zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50 zeta19wzjdtah4kl2vh77jks68cyy5gpjyurqltys99


# Execute all messages from executer . At this time zeta and mario keys can be offline
zetacored tx authz exec gas_price_voter_zeta.json --from executer --gas=auto --gas-prices=0.1azeta --gas-adjustment=10 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx authz exec nonce_voter_zeta.json --from executer --gas=auto --gas-prices=0.1azeta --gas-adjustment=10 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx authz exec mario_inbound_vote.json --from executer --gas=auto --gas-prices=0.1azeta --gas-adjustment=10 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx authz exec zeta_inbound_vote.json --from executer --gas=auto --gas-prices=0.1azeta --gas-adjustment=10 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx authz exec zeta_outbound_vote.json --from executer --gas=auto --gas-prices=0.1azeta --gas-adjustment=10 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx authz exec mario_outbound_vote.json --from executer --gas=auto --gas-prices=0.1azeta --gas-adjustment=10 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block





