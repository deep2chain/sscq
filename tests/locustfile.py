from locust import HttpLocust, TaskSet, task

class WebsiteTasks(TaskSet):
    # def on_start(self):
    #     self.client.post("/login", {
    #         "username": "test_user",
    #         "password": ""
    #     })
    
    @task
    def listaccounts(self):
        self.client.get("/accounts/list")
        
    @task
    def blocklastest(self):
        self.client.get("/blocks/latest")
    
    # @task
    # def blockinfo(self):
    #     self.client.get("/blocks/")
        
    @task
    def nodeinfo(self):
        self.client.get("/node_info")
        
    # @task
    # def validatorinfo(self):
    #     self.client.get("/validatorsets/")    

class WebsiteUser(HttpLocust):
    task_set = WebsiteTasks
    min_wait = 5000
    max_wait = 15000
