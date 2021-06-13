pub mod app;
pub mod routes;

fn main() {
    wasm_logger::init(wasm_logger::Config::new(log::Level::Trace));
    yew::start_app::<app::App>();
}
