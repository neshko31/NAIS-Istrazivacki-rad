package com.example.stakeholders_service.Repository;

import com.example.stakeholders_service.Model.Role;
import com.example.stakeholders_service.Model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByUserId(Long userId);
    List<User> findByRoleNot(Role role);
    boolean existsByUsername(String username);
    boolean existsByUserId(Long userId);
    List<User> findByRole(Role role);
}
