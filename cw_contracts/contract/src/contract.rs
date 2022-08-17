use cosmwasm_std::{entry_point, to_binary};
use cosmwasm_std::{Deps, DepsMut, Env, MessageInfo};
use cosmwasm_std::{QueryResponse, Response, StdError, StdResult};

use schemars::JsonSchema;
use thiserror::Error;

use serde::{Deserialize, Serialize};

use zeta_std::{OutTxTrackerAllResponse,ZetaCoreMsg,ZetaCoreQuery};

#[derive(Error, Debug)]
pub enum WatcherError {
    #[error("{0}")]
    Std(#[from] StdError),
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Instantiate
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)] //JsonSchema removed
pub struct InstantiateMsg {}

#[entry_point]
pub fn instantiate(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, WatcherError> {
    Ok(Response::default())
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Execute
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ExecuteMsg {
    AddToWatchList {
        chain: String,
        nonce: u32,
        tx_hash: String,
    }
}

#[entry_point]
pub fn execute(
    _deps: DepsMut<ZetaCoreQuery>,
    _env: Env,
    _info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response<ZetaCoreMsg>, WatcherError> {
    match msg {
        ExecuteMsg::AddToWatchList { chain,nonce,tx_hash } => {
            let add_watchlist_msg = ZetaCoreMsg::AddToWatchList {
                chain,
                nonce,
                tx_hash,
            };

            Ok(Response::new()
                .add_attribute("action", "add_watchlist")
                .add_message(add_watchlist_msg))
        }
    }
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Query
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum QueryMsg {
    OutTxTrackerAll {},
}

#[entry_point]
pub fn query(deps: Deps<ZetaCoreQuery>, _env: Env, msg: QueryMsg) -> StdResult<QueryResponse> {
    match msg {
        QueryMsg::OutTxTrackerAll {} => to_binary(&query_watchlist(deps)?),
    }
}

fn query_watchlist(deps: Deps<ZetaCoreQuery>) -> StdResult<OutTxTrackerAllResponse> {
    let req = ZetaCoreQuery::OutTxTrackerAll{}.into();
    deps.querier.query(&req)
}
