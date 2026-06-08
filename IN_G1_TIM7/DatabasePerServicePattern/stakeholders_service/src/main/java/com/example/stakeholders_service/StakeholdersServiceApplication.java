package com.example.stakeholders_service;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.client.discovery.EnableDiscoveryClient;

@SpringBootApplication
@EnableDiscoveryClient
public class StakeholdersServiceApplication {

	public static void main(String[] args) {
		SpringApplication.run(StakeholdersServiceApplication.class, args);
	}

}
