use cosmwasm_std::{CosmosMsg,CustomQuery};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ZetaCoreMsg {
    AddToWatchList {
        chain: String,
        nonce: u32,
        tx_hash: String,
    },
}


impl cosmwasm_std::CustomMsg for ZetaCoreMsg{}


impl From<ZetaCoreMsg> for CosmosMsg<ZetaCoreMsg> {
    fn from(original: ZetaCoreMsg) -> Self {
        CosmosMsg::Custom(original)
    }
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ZetaCoreQuery {
    OutTxTrackerAll {},
}

impl CustomQuery for ZetaCoreQuery {}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub struct OutTxTrackerAllResponse {
    pub out_tx_tracker : Vec<OutTxTracker>
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub struct OutTxTracker {
    pub index   : String ,
    pub chain   : String ,
    pub nonce   : String ,
    pub hashlist: Vec<TxHashList>,
}
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub struct TxHashList {
    pub txhash :String ,
    pub singer :String
}
