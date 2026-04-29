-- Insert sample students
INSERT INTO students (id, name, last_name, age, created_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Juan', 'García', 21, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440002', 'María', 'López', 22, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440003', 'Carlos', 'Rodríguez', 20, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440004', 'Ana', 'Martínez', 23, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440005', 'Luis', 'Sánchez', 21, CURRENT_TIMESTAMP);

-- Insert sample courses
INSERT INTO courses (id, name, description, credits, capacity, created_at) VALUES
('550e8400-e29b-41d4-a716-446655440101', 'Introducción a Go', 'Learn the basics of Go programming language', 4, 30, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440102', 'PostgreSQL Avanzado', 'Advanced PostgreSQL design and optimization', 3, 25, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440103', 'DevOps con Terraform', 'Infrastructure as Code using Terraform on AWS', 4, 20, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440104', 'RESTful API Design', 'Building scalable RESTful APIs', 3, 28, CURRENT_TIMESTAMP);

-- Insert sample enrollments
INSERT INTO enrollments (id, student_id, course_id, total_amount, created_at) VALUES
('550e8400-e29b-41d4-a716-446655440201', '550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440101', 450.00, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440202', '550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440102', 350.00, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440203', '550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440103', 500.00, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440204', '550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440102', 350.00, CURRENT_TIMESTAMP),
('550e8400-e29b-41d4-a716-446655440205', '550e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440104', 400.00, CURRENT_TIMESTAMP);

-- Query: List all students
SELECT id, name, last_name, age, created_at FROM students;

-- Query: List all courses with enrollment count
SELECT 
    c.id,
    c.name,
    c.description,
    c.credits,
    c.capacity,
    COUNT(e.id) AS enrolled_count
FROM courses c
LEFT JOIN enrollments e ON c.id = e.course_id
GROUP BY c.id, c.name, c.description, c.credits, c.capacity;

-- Query: List all enrollments with student and course details
SELECT 
    e.id,
    s.name AS student_name,
    s.last_name AS student_last_name,
    c.name AS course_name,
    e.total_amount,
    e.created_at
FROM enrollments e
JOIN students s ON e.student_id = s.id
JOIN courses c ON e.course_id = c.id;

-- Query: Count enrollments per student
SELECT 
    s.id,
    s.name,
    s.last_name,
    COUNT(e.id) AS enrollment_count,
    COALESCE(SUM(e.total_amount), 0) AS total_paid
FROM students s
LEFT JOIN enrollments e ON s.id = e.student_id
GROUP BY s.id, s.name, s.last_name;
