pub mod init {
    use crate::handlers::handlers;

    pub async fn new_service() {
        tokio::spawn(
            handlers::init_http()
        ).await;
    }
}
