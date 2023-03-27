use axum::{
    Router,
    routing::get,
    response::{Headers, IntoResponse},
};

pub mod handlers {
    async fn init_http() {
        let app = Router::new().route("/volumes", get(list_volumes));
    }

    async fn list_volumes() -> impl IntoResponse {}
}
