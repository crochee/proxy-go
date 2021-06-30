pub mod home;
pub mod not_found;

use yew_router::prelude::*;

#[derive(Switch, Debug, PartialEq)]
pub enum AppRoute {
    #[to = "/"]
    Home,
    #[to = "/404"]
    NotFound,
}
