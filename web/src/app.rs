use yew::{html, Component, ComponentLink, Html, ShouldRender};
use yew_router::prelude::*;

use crate::routes::{home::Home, not_found::NotFound, AppRoute};

pub enum Msg {
    Route(Route),
}

pub struct App {
    current_route: Option<AppRoute>,
    link: ComponentLink<Self>,
}

impl Component for App {
    type Message = Msg;
    type Properties = ();

    fn create(_props: Self::Properties, link: ComponentLink<Self>) -> Self {
        let route_service: RouteService = RouteService::new();
        let route = route_service.get_route();
        App {
            current_route: AppRoute::switch(route),
            link,
        }
    }

    fn update(&mut self, msg: Self::Message) -> ShouldRender {
        match msg {
            Msg::Route(route) => {
                self.current_route = AppRoute::switch(route);
                true
            }
        }
    }

    fn change(&mut self, _props: Self::Properties) -> ShouldRender {
        false
    }

    fn view(&self) -> Html {
        html! {
            <>
                <main>
                {
                    // Routes to render sub components
                    if let Some(route) = &self.current_route {
                        match route {
                            AppRoute::Home=>html!{<Home/>},
                            AppRoute::NotFound=> html! {<NotFound/>},
                        }
                    } else {
                        html!{<NotFound/>}
                    }
                }
                </main>
                <footer class="footer">
                    <div class="content has-text-centered">
                        { "Powered by " }
                        <a href="https://yew.rs">{ "Yew" }</a>
                        { " using " }
                        <a href="https://bulma.io">{ "Bulma" }</a>
                        { " and images from " }
                        <a href="https://unsplash.com">{ "Unsplash" }</a>
                    </div>
                </footer>
            </>
        }
    }
}
