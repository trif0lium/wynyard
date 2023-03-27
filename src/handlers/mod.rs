pub mod handlers {
    use axum::{
        Router,
        routing::get,
        response::IntoResponse
    };

    pub async fn init_http() {
        let app = Router::new().route("/volumes", get(list_volumes));
    }

    async fn list_volumes() -> impl IntoResponse {}
}
