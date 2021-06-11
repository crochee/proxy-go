use yew::prelude::*;
use yew_router::prelude::*;

use page::{home::Home, not_found::NotFound};

mod page;

#[derive(Switch, PartialEq, Clone, Debug)]
pub enum RoutePath {
    #[to = "/"]
    Home,
    #[not_found]
    #[to = "/404"]
    NotFound,
}

fn switch_route(routes: &RoutePath) -> Html {
    match routes {
        RoutePath::Home => html! { <Home /> },
        RoutePath::NotFound => html! { <NotFound /> },
    }
}

pub enum Msg {
    ToggleNavbar,
}

pub struct Model {
    link: ComponentLink<Self>,
    navbar_active: bool,
}

impl Component for Model {
    type Message = Msg;
    type Properties = ();

    fn create(_props: Self::Properties, link: ComponentLink<Self>) -> Self {
        Self {
            link,
            navbar_active: false,
        }
    }

    fn update(&mut self, msg: Self::Message) -> ShouldRender {
        match msg {
            Msg::ToggleNavbar => {
                self.navbar_active = !self.navbar_active;
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
                { self.view_nav() }

                <main>
                    <Router<Route> render=Router::render(switch_route) />
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

impl Model {
    fn view_nav(&self) -> Html {
        let Self {
            ref link,
            navbar_active,
            ..
        } = *self;

        let active_class = if navbar_active { "is-active" } else { "" };

        html! {
            <nav class="navbar is-primary" role="navigation" aria-label="main navigation">
                <div class="navbar-brand">
                    <h1 class="navbar-item is-size-3">{ "Yew Blog" }</h1>

                    <a role="button"
                        class=classes!("navbar-burger", "burger", active_class)
                        aria-label="menu" aria-expanded="false"
                        onclick=link.callback(|_| Msg::ToggleNavbar)
                    >
                        <span aria-hidden="true"></span>
                        <span aria-hidden="true"></span>
                        <span aria-hidden="true"></span>
                    </a>
                </div>
                <div class=classes!("navbar-menu", active_class)>
                    <div class="navbar-start">
                        <Link<Route> classes=classes!("navbar-item") route=Route::Home>
                            { "Home" }
                        </Link<Route>>
                    </div>
                </div>
            </nav>
        }
    }
}


fn main() {
    yew::start_app::<Model>();
}