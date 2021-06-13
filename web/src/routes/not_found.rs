use yew::{html, Component, ComponentLink, Html, ShouldRender};

pub struct NotFound {}

impl Component for NotFound {
    type Message = ();
    type Properties = ();

    fn create(_props: Self::Properties, _link: ComponentLink<Self>) -> Self {
        Self { /* fields */ }
    }

    fn update(&mut self, _msg: Self::Message) -> ShouldRender {
        unimplemented!()
    }

    fn change(&mut self, _: <Self as yew::Component>::Properties) -> bool {
        false
    }

    fn view(&self) -> Html {
        html! {
            <div>
                {"404 Not Found"}
            </div>
        }
    }
}
