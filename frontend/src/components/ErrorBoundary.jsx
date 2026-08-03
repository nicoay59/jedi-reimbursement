import React from "react";

export default class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      hasError: false,
      message: ""
    };
  }

  static getDerivedStateFromError(error) {
    return {
      hasError: true,
      message:
        error?.message ||
        "Terjadi kesalahan yang tidak terduga."
    };
  }

  componentDidCatch(error, info) {
    console.error("Application error:", error, info);
  }

  handleReload = () => {
    window.location.reload();
  };

  handleHome = () => {
    window.location.assign("/");
  };

  render() {
    if (!this.state.hasError) {
      return this.props.children;
    }

    return (
      <main className="min-vh-100 d-grid place-items-center p-4">
        <div className="card app-card error-boundary-card">
          <div className="card-body p-5 text-center">
            <div className="error-code">500</div>
            <h1 className="h3">Aplikasi mengalami kesalahan</h1>
            <p className="text-secondary">
              {this.state.message}
            </p>

            <div className="d-flex flex-column flex-sm-row justify-content-center gap-2 mt-4">
              <button
                type="button"
                className="btn btn-primary"
                onClick={this.handleReload}
              >
                Muat ulang
              </button>

              <button
                type="button"
                className="btn btn-outline-secondary"
                onClick={this.handleHome}
              >
                Kembali ke beranda
              </button>
            </div>
          </div>
        </div>
      </main>
    );
  }
}
