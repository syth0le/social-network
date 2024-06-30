box.cfg {
    listen = "3301"
}

box.once(
        "schema",
        function()
            box.schema.space.create('users', { if_not_exists = true })

            box.space.users:format({
                { name = 'id', type = 'unsigned' },
                { name = 'first_name', type = 'string' },
                { name = 'second_name', type = 'string' },
                { name = 'username', type = 'string' },
                { name = 'hashed_password', type = 'string' },
                { name = 'sex', type = 'string' },
                { name = 'biography', type = 'string' },
                { name = 'city', type = 'string' },
            })

            box.space.users:create_index('primary', { type = "TREE", unique = true, parts = { 1, 'unsigned' }, if_not_exists = true })
            box.space.users:create_index('first_second_name_idx', { type = 'TREE', unique = false, parts = {2, 'string', 3, 'string' }, if_not_exists = true })
            box.space.users:create_index('first_name_idx', { type = 'TREE', unique = false, parts = { 2, 'string' }, if_not_exists = true })
            box.space.users:create_index('second_name_idx', { type = 'TREE', unique = false, parts = { 3, 'string' }, if_not_exists = true })

            box.space.users:insert({1, "first_name", "second_name", "username", "password", "male", "bio", "Berlin"})
        end
)

-- procedure for search by first name prefix AND second name prefix
-- Param: prefix_first_name - prefix for searching first name by like '%first_name'
-- Param: prefix_second_name - prefix for searching second name by like '%second_name'
-- Param: size - max count of entries in response
-- Param: offset - offset from data start position
function search_by_first_second_name_with_offset(prefix_first_name, prefix_second_name, size, offset)
    local count = 0
    local step = 0
    local result = {}
    for _, tuple in box.space.users.index.first_second_name_idx:pairs(prefix_first_name, { iterator = 'GE' }) do
        if string.startswith(tuple[2], prefix_first_name, 1, -1) and string.startswith(tuple[3], prefix_second_name, 1, -1) then
            if step < offset then
                step = step + 1
            else
                table.insert(result, tuple)
                count = count + 1
                if count >= size then
                    return result
                end
            end
        end
    end
    return result
end

function search_by_first_second_name_with_size(prefix_first_name, prefix_second_name, size)
    local count = 0
    local result = {}
    for _, tuple in box.space.users.index.first_second_name_idx:pairs(prefix_first_name, { iterator = 'GE' }) do
        if string.startswith(tuple[2], prefix_first_name, 1, -1) and string.startswith(tuple[3], prefix_second_name, 1, -1) then
            table.insert(result, tuple)
            count = count + 1
            if count >= size then
                return result
            end
        end
    end
    return result
end

function search_by_first_second_name(prefix_first_name, prefix_second_name)
    local result = {}
    for _, tuple in box.space.users.index.first_second_name_idx:pairs(prefix_first_name, { iterator = 'GE' }) do
        if string.startswith(tuple[2], prefix_first_name, 1, -1) and string.startswith(tuple[3], prefix_second_name, 1, -1) then
            table.insert(result, tuple)
        end
    end
    return result
end

-- procedure for search by first name prefix
-- Param: prefix - prefix for searching first name by like '%first_name'
function search_by_first_name(prefix)
    local result = {}
    for _, tuple in box.space.users.index.first_name_idx:pairs({ prefix }, { iterator = 'GE' }) do
        if string.startswith(tuple[2], prefix, 1, -1) then
            table.insert(result, tuple)
        end
    end
    return result
end

-- procedure for search by second name prefix
-- Param: prefix - prefix for searching first name by like '%second_name'
function search_by_second_name(prefix)
    local result = {}
    for _, tuple in box.space.users.index.second_name_idx:pairs({ prefix }, { iterator = 'GE' }) do
        if string.startswith(tuple[3], prefix, 1, -1) then
            table.insert(result, tuple)
        end
    end
    return result
end