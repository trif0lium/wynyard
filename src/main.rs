mod init;

#[tokio::main]
async fn main() {
    init::init::new_service().await
}
