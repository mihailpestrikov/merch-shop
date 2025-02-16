from locust import HttpUser, task, between

class MyUser(HttpUser):
    wait_time = between(0.05, 0.2)

    @task
    def auth_and_get_info(self):
        auth_response = self.client.post(
            "/api/auth", json={"username": "testuser", "password": "password"})

        if auth_response.status_code == 200:
            token = auth_response.json().get("token")

            headers = {"Authorization": f"Bearer {token}"}
            self.client.get("/api/info", headers=headers)

if __name__ == "__main__":
    import os
    os.system("locust -f locustfile.py --run-time 3m --host http://localhost:8080 "
              "--web-host=127.0.0.1 --web-port=8089")
