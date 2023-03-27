mod init;
mod controllers;
mod handlers;
mod gateways;

#[tokio::main]
async fn main() {
    init::init::new_service().await
}
