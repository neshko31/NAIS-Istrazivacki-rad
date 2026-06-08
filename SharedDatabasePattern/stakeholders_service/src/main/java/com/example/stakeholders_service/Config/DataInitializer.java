package com.example.stakeholders_service.Config;

import com.example.stakeholders_service.Model.Role;
import com.example.stakeholders_service.Model.User;
import com.example.stakeholders_service.Repository.UserRepository;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.stereotype.Component;

/**
 * Seeds profile data into the shared users table.
 *
 * NOTE (antipattern demo): this service writes directly to the same
 * users table that auth-service owns. In a proper architecture, each
 * service would own its own table and communicate via API.
 *
 * Passwords are intentionally NOT set here — that is auth-service's
 * responsibility. The password column in users belongs to auth-service.
 */
@Component
public class DataInitializer implements ApplicationRunner {

    private final UserRepository userRepository;

    public DataInitializer(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    @Override
    public void run(ApplicationArguments args) {
        insertIfNotExists(1L, "marija",  "marija@example.com",  Role.TOURIST, "Marija",  "Todorovic");
        insertIfNotExists(2L, "nenad",   "nenad@example.com",   Role.ADMIN,   "Nenad",   "Lukic");
        insertIfNotExists(3L, "teodora", "teodora@example.com", Role.TOURIST, "Teodora", "Pesic");
        insertIfNotExists(4L, "srdjan",  "srdjan@example.com",  Role.GUIDE,   "Srđan",   "Sancanin");
    }

    private void insertIfNotExists(Long registryId, String username, String email,
                                   Role role, String firstName, String lastName) {
        if (!userRepository.existsByUsername(username)) {
            User user = new User();
            user.setUserId(registryId);
            user.setUsername(username);
            user.setEmail(email);
            user.setRole(role);
            user.setFirstName(firstName);
            user.setLastName(lastName);
            user.setBlocked(false);
            userRepository.save(user);
            System.out.println("[DataInitializer] Created user profile: " + username + " (registryId=" + registryId + ")");
        }
    }
}
